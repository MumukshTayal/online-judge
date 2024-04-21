import React, { useEffect, useState } from 'react';
import { useNavigate, Route, useParams } from 'react-router-dom';
import Timer from './timer.tsx';
import { CardTitle, CardHeader, CardContent, CardFooter, Card } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import Navbar from "./navbar.tsx"

export default function ContestView() {
  const { contestId } = useParams();
  const [contest, setContest] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchContest = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/get_contest_details/${contestId}`);
        if (!response.ok) {
          throw new Error('Failed to fetch contest');
        }
        const data = await response.json();
        setContest(data);
      } catch (error) {
        console.error(error);
      }
    };

    fetchContest();
  }, [contestId]);

  const handleProblemClick = (problemId) => {
    navigate(`/contest/${contestId}/problem/${problemId}`);
  };

  if (!contest) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <Navbar />
      <main className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="grid gap-4 w-full max-w-3xl p-4 rounded-lg border dark:border-gray-800">
          <div className="flex items-center gap-4">
            <h1 className="text-3xl font-bold tracking-tighter">{contest.contest.contest_title}</h1>
            <div className="ml-auto flex items-center gap-2">
              <div className="flex items-center gap-0.5">
                <ClockIcon className="w-4 h-4" />
                <span className="font-semibold">Time remaining:</span>
              </div>
              <div className="flex items-center gap-0.5">
                <Timer endTime={contest.contest.end_time} />
              </div>
            </div>
          </div>
          <div className="grid gap-4">
            {contest.contest_problems.map(problem => (
              <Card key={problem.problem_id}>
                <CardHeader className="flex items-center space-x-2">
                  <CardTitle>{problem.problem_title}</CardTitle>
                  <Badge variant={problem.status === 'Unattempted' ? 'neutral' : problem.status === 'Accepted' ? 'positive' : 'negative'}>
                    {problem.status}
                  </Badge>
                </CardHeader>
                <CardContent>
                  <p className="text-sm leading-none">
                    Description: {problem.problem_description}
                  </p>
                </CardContent>
                <CardFooter>
                  {/* Add onClick event to navigate to ProblemView with problem_id */}
                  <Button size="sm" onClick={() => handleProblemClick(problem.problem_id)}>View</Button>
                </CardFooter>
              </Card>
            ))}
          </div>
          <div className="flex flex-col gap-1 min-[400px]:flex-row justify-center">
            <a
              className="inline-flex h-10 items-center justify-center rounded-md border border-gray-200 border-gray-200 bg-white px-8 text-sm font-medium shadow-sm gap-1 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:border-gray-800 dark:border-gray-800 dark:bg-gray-950 dark:hover:bg-gray-950 dark:hover:text-gray-50 dark:focus-visible:ring-gray-300"
              href={`/contest/${contestId}/submissions`}
            >
              Submissions
            </a>
            <a
              className="inline-flex h-10 items-center justify-center rounded-md bg-gray-900 px-8 text-sm font-medium text-gray-50 shadow gap-1 transition-colors hover:bg-gray-900/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-gray-950 disabled:pointer-events-none disabled:opacity-50 dark:bg-gray-50 dark:text-gray-900 dark:hover:bg-gray-50/90 dark:focus-visible:ring-gray-300"
              href={`/contest/${contestId}/leaderboard`}
            >
              Leaderboard
            </a>
          </div>
        </div>
      </main>
    </div>
  );
}

function CodeIcon(props) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <polyline points="16 18 22 12 16 6" />
      <polyline points="8 6 2 12 8 18" />
    </svg>
  );
}

function ClockIcon(props) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <circle cx="12" cy="12" r="10" />
      <polyline points="12 6 12 12 16 14" />
    </svg>
  )
}
