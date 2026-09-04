import { Navigate, createBrowserRouter } from 'react-router-dom'

import { AuthPage } from '../pages/AuthPage'
import { ProtectedRoute } from '../features/auth/ProtectedRoute'
import { LibraryPage } from '../pages/LibraryPage'
import { NotFoundPage } from '../pages/NotFoundPage'
import { ReaderPage } from '../pages/ReaderPage'
import { VocabularyPage } from '../pages/VocabularyPage'
import { AppLayout } from '../shared/layout/AppLayout'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Navigate to="/library" replace />,
  },
  {
    path: '/login',
    element: <AuthPage mode="login" />,
  },
  {
    path: '/register',
    element: <AuthPage mode="register" />,
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <AppLayout />,
        children: [
          {
            path: '/library',
            element: <LibraryPage kind="my-library" />,
          },
          {
            path: '/catalog',
            element: <LibraryPage kind="catalog" />,
          },
          {
            path: '/vocabulary',
            element: <VocabularyPage />,
          },
        ],
      },
      {
        path: '/reader/:bookId',
        element: <ReaderPage />,
      },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
