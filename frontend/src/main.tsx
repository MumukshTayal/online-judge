import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import HomePage from './components/home-page.tsx'
import CreateContest from './components/create-contest.tsx'
import ContestList from './components/contest-list.tsx'
import SubmissionsPage from './components/submissions_page.tsx'
import AddProblemPage from './components/add-problem-page.tsx'
import AddTestCasePage from './components/add-test-case-page.tsx'
import ContestView from './components/contest-view.tsx'
import './index.css'

import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom"; 

const router = createBrowserRouter([
  {
    path: "/home", 
    element: <HomePage />,
  },
  {
    path: "/", 
    element: <App />,
  },
  {
    path: "/contest-list",
    element: <ContestList />
  },
  {
    path: "/contest/:contestId", 
    element: <ContestView />
  },
  {
    path: "/create-contest",
    element: <CreateContest />
  },
  {
    path: "/submissions",
    element: <SubmissionsPage />
  },
  {
    path: "/add-problem",
    element: <AddProblemPage />
  },
  {
    path: "/add-testcase",
    element: <AddTestCasePage />
  },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
