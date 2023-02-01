import React from 'react';
import ReactDOM from 'react-dom/client';
import { createMemoryRouter, RouterProvider } from 'react-router-dom';
import { App } from './app';
import Files from "./menu/files";
import Peers from './menu/peers';
import Peer from './menu/peers/Peer';
import './style.css';

const router = createMemoryRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { path: '/peers', element: <Peers />, children: [{ path: '/peers/:id', element: <Peer /> }] },
      { path: '/files', element: <Files /> },
    ]
  }
]);

ReactDOM.createRoot(document.getElementById('app') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
