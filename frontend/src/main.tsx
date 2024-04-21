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
import ProblemView from './components/problem-view.tsx'

import './index.css'

import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom"; 
import LeaderboardPage from './components/leaderboard.tsx'
import ProblemsPage from './components/problems-page.tsx'
import AllSubmissionsPage from './components/all_submissions.tsx'

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
    path: "/contest/:contestId/problem/:problemId", 
    element: <ProblemView />
  },
  {
    path: "/create-contest",
    element: <CreateContest />
  },
  {
    path: "/submissions",
    element: <AllSubmissionsPage/>
  },
  {
    path: "/add-problem",
    element: <AddProblemPage />
  },
  {
    path: "/problems",
    element: <ProblemsPage />
  },
  {
    path: "/add-testcase",
    element: <AddTestCasePage />
  },
  {
    path: "contest/:contestId/leaderboard",
    element: <LeaderboardPage/>
  },
  {
    path: "contest/:contestId/submissions",
    element: <SubmissionsPage/>
  }
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
