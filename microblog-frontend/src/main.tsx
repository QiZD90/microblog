import React from 'react';
import ReactDOM from 'react-dom/client';
import * as Router from 'react-router-dom';
import { ChakraProvider } from '@chakra-ui/react';
import MainPage from './pages/MainPage';

const router = Router.createBrowserRouter([
  {
    path: '/',
    element: <MainPage />,
  },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ChakraProvider>
      <Router.RouterProvider router={router} />
    </ChakraProvider>
  </React.StrictMode>
);
