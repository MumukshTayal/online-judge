import { useState, useEffect } from 'react';
import { CardContent, Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useNavigate, Route } from 'react-router-dom';
import ContestView from './contest-view.tsx';
import { Link } from 'react-router-dom';
import Navbar from "./navbar.tsx"

export default function ContestList() {
  const navigate = useNavigate();
  const [participatingContests, setParticipatingContests] = useState([]);
  const [createdContests, setCreatedContests] = useState([]);

  useEffect(() => {
    fetchContests();
  }, []);

  const fetchContests = async () => {
    try {
      const jwtToken = localStorage.getItem('jwtToken');
  
      const response = await fetch("http://localhost:8080/api/get_contest_list", {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${jwtToken}`,
          'Content-Type': 'application/json'
        },
      });
  
      if (!response.ok) {
        throw new Error('Failed to fetch contests');
      }
  
      const data = await response.json();
      setParticipatingContests(data.participating_contests);
      setCreatedContests(data.created_contests);
    } catch (error) {
      console.error(error);
    }
  };
  

  const handleContestClick = (contest) => {
    navigate(`/contest/${contest.contest_id}`);
  };

  return (
    <div className="flex flex-col h-screen">
      <Navbar />
      <main className="flex-1 overflow-y-auto">
        <section className="container py-6 space-y-6 text-gray-900 md:space-y-8 dark:text-gray-50">
          <div className="space-y-2">
            <h1 className="text-3xl font-semibold tracking-tighter sm:text-4xl md:text-5xl">Contests</h1>
            <br></br>
            <p className="max-w text-gray-500 md:text-base/relaxed dark:text-gray-400">
            Expand our contest collection by adding new challenges. Contribute to our platform's growth and provide engaging competitions for participants.
            </p>
            <Link to="/create-contest">
              <Button variant="primary" className="w-32 bg-black text-white mt-4">
                Add Contest
              </Button>
            </Link>
          </div>
          
          {/* Display participating contests */}
          <div className="grid grid-cols-1 gap-6">
            <h2 className="text-2xl font-semibold">Participate</h2>
            <p className="max-w-prose text-gray-500 md:text-base/relaxed dark:text-gray-400">
              These are your available contests. Participate and improve your skills.
            </p>
            {participatingContests && participatingContests.length > 0 ? ( 
              participatingContests.map(contest => (
                <Card key={contest.contest_id} onClick={() => handleContestClick(contest)} className="cursor-pointer">
                  <CardContent className="p-4 md:p-6">
                    <div className="space-y-2">
                      <h3 className="text-lg font-semibold">{contest.contest_title}</h3>
                      <p className="text-sm text-gray-500 dark:text-gray-400">{contest.contest_description}</p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">{contest.is_public}</p>
                      <div className="flex items-center space-x-2 text-sm">
                        <ClockIcon className="w-4 h-4 text-gray-500" />
                        <time dateTime={contest.start_time}>{new Date(contest.start_time).toLocaleString()}</time>
                      </div>
                      <div className="flex items-center space-x-2 text-sm">
                        <ClockIcon className="w-4 h-4" />
                        <time dateTime={contest.end_time}>{new Date(contest.end_time).toLocaleString()}</time>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <p>No participating contests available.</p>
            )}
          </div>
          
          {/* Display created contests */}
          <div className="grid grid-cols-1 gap-6">
            <h2 className="text-2xl font-semibold">Created Contests</h2>
            {createdContests && createdContests.length > 0 ? ( 
              createdContests.map(contest => (
                <Card key={contest.contest_id} onClick={() => handleContestClick(contest)} className="cursor-pointer">
                  <CardContent className="p-4 md:p-6">
                    <div className="space-y-2">
                      <h3 className="text-lg font-semibold">{contest.contest_title}</h3>
                      <p className="text-sm text-gray-500 dark:text-gray-400">{contest.contest_description}</p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">{contest.is_public}</p>
                      <div className="flex items-center space-x-2 text-sm">
                        <ClockIcon className="w-4 h-4 text-gray-500" />
                        <time dateTime={contest.start_time}>{new Date(contest.start_time).toLocaleString()}</time>
                      </div>
                      <div className="flex items-center space-x-2 text-sm">
                        <ClockIcon className="w-4 h-4" />
                        <time dateTime={contest.end_time}>{new Date(contest.end_time).toLocaleString()}</time>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <p>No created contests available.</p>
            )}
          </div>
        </section>
      </main>
    </div>
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

// Add the Route component to handle the ContestView component
<Route path="/contest/:contestId" element={<ContestView />} />
