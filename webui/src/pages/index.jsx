import React, { useEffect } from 'react';

import { BrowserRouter as Router, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { WithLoginConfigContext } from '../lib/hooks/conf';

// pages
import RepositoriesPage from './repositories';
import { RepositoryPageLayout } from '../lib/components/repository/layout.jsx';
import RepositoryObjectsPage from './repositories/repository/objects';
import RepositoryObjectsViewPage from './repositories/repository/objectViewer';

import RepositoryCommitsPage from './repositories/repository/commits';
import RepositoryCommitPage from './repositories/repository/commits/commit';
import RepositoryBranchesPage from './repositories/repository/branches';
import RepositoryRevertPage from './repositories/repository/branches/revert';
import RepositoryTagsPage from './repositories/repository/tags';
import RepositoryComparePage from './repositories/repository/compare';
import Layout from '../lib/components/layout';
import LoginPage from './auth/login';
import { WithAppContext } from '../lib/hooks/appContext';
import { AuthProvider } from '../lib/auth/authContext';
import RequiresAuth from '../lib/components/requiresAuth';

export const IndexPage = () => {
    return (
        <Router>
            <AuthProvider>
                <WithAppContext>
                    <WithLoginConfigContext>
                        <Routes>
                            <Route element={<RequiresAuth />}>
                                <Route index element={<Navigate to="/repositories" />} />
                                <Route path="repositories" element={<Layout />}>
                                    <Route index element={<RepositoriesPage />} />
                                    <Route path=":repoId" element={<RepositoryPageLayout />}>
                                        <Route path="objects" element={<RepositoryObjectsPage />} />
                                        <Route path="object" element={<RepositoryObjectsViewPage />} />
                                        <Route path="commits">
                                            <Route index element={<RepositoryCommitsPage />} />
                                            <Route path=":commitId" element={<RepositoryCommitPage />} />
                                        </Route>
                                        <Route path="branches">
                                            <Route index element={<RepositoryBranchesPage />} />
                                            <Route path=":branchId/revert" element={<RepositoryRevertPage />} />
                                        </Route>
                                        <Route path="tags" element={<RepositoryTagsPage />} />
                                        <Route path="compare/*" element={<RepositoryComparePage />} />
                                        <Route index element={<Navigate to="objects" />} />
                                    </Route>
                                </Route>
                                <Route path="*" element={<Navigate to="/repositories" replace />} />
                            </Route>
                            <Route path="auth" element={<Layout />}>
                                <Route path="login" element={<LoginPage />} />
                            </Route>
                        </Routes>
                    </WithLoginConfigContext>
                </WithAppContext>
            </AuthProvider>
        </Router>
    );
};
