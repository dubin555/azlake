import * as duckdb from '@duckdb/duckdb-wasm';
import * as arrow from 'apache-arrow';
import { AsyncDuckDB, AsyncDuckDBConnection } from '@duckdb/duckdb-wasm';
import duckdb_wasm from '@duckdb/duckdb-wasm/dist/duckdb-mvp.wasm?url';
import mvp_worker from '@duckdb/duckdb-wasm/dist/duckdb-browser-mvp.worker.js?url';
import duckdb_wasm_eh from '@duckdb/duckdb-wasm/dist/duckdb-eh.wasm?url';
import eh_worker from '@duckdb/duckdb-wasm/dist/duckdb-browser-eh.worker.js?url';

const MANUAL_BUNDLES: duckdb.DuckDBBundles = {
    mvp: {
        mainModule: duckdb_wasm,
        mainWorker: mvp_worker,
    },
    eh: {
        mainModule: duckdb_wasm_eh,
        mainWorker: eh_worker,
    },
};

let _db: AsyncDuckDB | null = null;

async function getDuckDB(): Promise<duckdb.AsyncDuckDB> {
    if (_db !== null) {
        return _db;
    }
    const bundle = await duckdb.selectBundle(MANUAL_BUNDLES);
    if (!bundle.mainWorker) {
        throw Error('could not initialize DuckDB');
    }
    const worker = new Worker(bundle.mainWorker);
    const logger = new duckdb.VoidLogger();
    const db = new duckdb.AsyncDuckDB(logger, worker);
    await db.instantiate(bundle.mainModule, bundle.pthreadWorker);
    const conn = await db.connect();
    await conn.close();
    _db = db;
    return _db;
}

// taken from @duckdb/duckdb-wasm/dist/types/src/bindings/tokens.d.ts
const DUCKDB_STRING_CONSTANT = 2;
const LAKEFS_URI_PATTERN = /^(['"]?)(lakefs:\/\/(.*))(['"])\s*$/;

const API_BASE = `${document.location.protocol}//${document.location.host}/api/v1`;

// Resolve a lakefs:// URI to a download URL.
// Tries SAS URL first (direct Azure access), falls back to REST API.
async function resolveObjectUrl(repo: string, ref: string, path: string): Promise<string> {
    try {
        const resp = await fetch(`${API_BASE}/repositories/${repo}/refs/${ref}/objects/sas?path=${encodeURIComponent(path)}`);
        if (resp.ok) {
            const data = await resp.json();
            if (data.sas_url) {
                return data.sas_url;
            }
        }
    } catch {
        // SAS not available, fall through
    }
    return `${API_BASE}/repositories/${repo}/refs/${ref}/objects?path=${encodeURIComponent(path)}`;
}

// Parse SQL to find lakefs:// URIs, resolve them, and return file info
async function extractFiles(conn: AsyncDuckDBConnection, sql: string): Promise<{ original: string; url: string; restUrl: string; localName: string }[]> {
    const tokenized = await conn.bindings.tokenize(sql);
    let prev = 0;
    const files: { original: string; url: string; restUrl: string; localName: string }[] = [];
    const seen = new Set<string>();

    for (let i = 0; i < tokenized.offsets.length; i++) {
        let currentToken = sql.length;
        if (i < tokenized.offsets.length - 1) {
            currentToken = tokenized.offsets[i + 1];
        }
        const part = sql.substring(prev, currentToken);
        prev = currentToken;

        if (tokenized.types[i] === DUCKDB_STRING_CONSTANT) {
            const matches = part.match(LAKEFS_URI_PATTERN);
            if (matches !== null) {
                const unescapedUri = matches[2].replace(/''/g, "'");
                if (seen.has(unescapedUri)) continue;
                seen.add(unescapedUri);

                const pathParts = matches[3].replace(/''/g, "'").split('/');
                const repo = pathParts[0];
                const ref = pathParts[1];
                const filePath = pathParts.slice(2).join('/');
                const fileName = pathParts[pathParts.length - 1];

                const restUrl = `${API_BASE}/repositories/${repo}/refs/${ref}/objects?path=${encodeURIComponent(filePath)}`;
                const url = await resolveObjectUrl(repo, ref, filePath);
                const localName = `/data/${files.length}_${fileName}`;
                files.push({ original: unescapedUri, url, restUrl, localName });
            }
        }
    }
    return files;
}

// Replace lakefs:// URIs in SQL with local DuckDB file names
function rewriteSQL(sql: string, files: { original: string; localName: string }[]): string {
    let rewritten = sql;
    for (const f of files) {
        rewritten = rewritten.replaceAll(`'${f.original}'`, `'${f.localName}'`);
        rewritten = rewritten.replaceAll(`"${f.original}"`, `"${f.localName}"`);
        const escaped = f.original.replace(/'/g, "''");
        rewritten = rewritten.replaceAll(`'${escaped}'`, `'${f.localName}'`);
    }
    return rewritten;
}

/* eslint-disable  @typescript-eslint/no-explicit-any */
export async function runDuckDBQuery(sql: string): Promise<arrow.Table<any>> {
    const db = await getDuckDB();
    let result: arrow.Table<any>;
    const conn = await db.connect();
    const registeredFiles: string[] = [];
    try {
        // Extract lakefs:// URIs and resolve to download URLs (SAS or REST API)
        const files = await extractFiles(conn, sql);

        // Fetch each file via JS fetch() and register as in-memory buffer in DuckDB.
        // This avoids needing httpfs extension or S3 gateway entirely.
        for (const f of files) {
            let resp: Response | null = null;
            // Try primary URL (SAS or REST API)
            try {
                resp = await fetch(f.url);
                if (!resp.ok) resp = null;
            } catch {
                resp = null;
            }
            // Fallback to REST API if SAS URL failed (e.g., CORS)
            if (!resp && f.restUrl !== f.url) {
                resp = await fetch(f.restUrl);
            }
            if (!resp || !resp.ok) {
                throw new Error(`Failed to fetch file data`);
            }
            const buffer = new Uint8Array(await resp.arrayBuffer());
            await db.registerFileBuffer(f.localName, buffer);
            registeredFiles.push(f.localName);
        }

        // Rewrite SQL to use local file names and execute
        const rewrittenSQL = rewriteSQL(sql, files);
        result = await conn.query(rewrittenSQL);

        // Cleanup registered files
        for (const name of registeredFiles) {
            await db.dropFile(name);
        }
    } finally {
        await conn.close();
    }
    return result;
}
