import { useState, useEffect } from 'react';
import { CardContent, Card } from "@/components/ui/card";
import { useNavigate } from 'react-router-dom';

export default function ContestList() {
  const navigate = useNavigate();
  const [contests, setContests] = useState([]);

  useEffect(() => {
    fetchContests();
  }, []);

  const fetchContests = async () => {
    try {
      const response = await fetch("/api/get_contest");
      if (!response.ok) {
        throw new Error('Failed to fetch contests');
      }
      const data = await response.json();
      setContests(data);
    } catch (error) {
      console.error(error);
    }
  };

  const handleCreateContestClick = () => {
    navigate("/create-contest");
  }; 

  return (
    <div className="flex flex-col h-screen">
      <header className="px-4 lg:px-6 h-14 flex items-center">
        <nav>
          <a className="flex items-center justify-center" href="/home">
            <CodeIcon className="h-6 w-6" />
            <span className="sr-only">Online Judge</span>
          </a> 
        </nav>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/contest-list">
            Contests
          </a>
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/add-problem">
            Add Problem 
          </a>
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/submissions">
            Submissions 
          </a>
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/add-testcase">
            Add Test Cases 
          </a>
        </nav>
      </header>
      <main className="flex-1 overflow-y-auto">
        <section className="container py-6 space-y-6 text-gray-900 md:space-y-8 dark:text-gray-50">
          <div className="space-y-2">
            <h1 className="text-3xl font-semibold tracking-tighter sm:text-4xl md:text-5xl">Contests</h1>
            <p className="max-w-prose text-gray-500 md:text-base/relaxed dark:text-gray-400">
              Participate in the latest contests and improve your skills.
            </p>
          </div>
          <div className="grid grid-cols-1 gap-6"> 
            <Card onClick={handleCreateContestClick} className="cursor-pointer">
              <CardContent className="p-4 md:p-6">
                {/* Create Contest button */}
                <button className="text-lg font-semibold text-gray-900 hover:text-gray-700 focus:outline-none">
                  Create Contest
                </button>
              </CardContent>
            </Card>
            {/* Display fetched contests */}
            {contests.map(contest => (
              <Card key={contest.id}>
                <CardContent className="p-4 md:p-6 cursor-pointer">
                  <div className="space-y-2">
                    <h3 className="text-lg font-semibold">{contest.name}</h3>
                    <p className="text-sm text-gray-500 dark:text-gray-400">{contest.description}</p>
                    <div className="flex items-center space-x-2 text-sm">
                      <ClockIcon className="w-4 h-4 text-gray-500" />
                      <time dateTime={contest.startTime}>{contest.startTime}</time>
                    </div>
                    <div className="flex items-center space-x-2 text-sm">
                      <ClockIcon className="w-4 h-4" />
                      <time dateTime={contest.endTime}>{contest.endTime}</time>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>
      </main>
    </div>
  )
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
  )
}
