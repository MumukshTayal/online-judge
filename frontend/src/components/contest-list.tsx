import { CardContent, Card } from "@/components/ui/card";
import { useNavigate } from 'react-router-dom';

export default function ContestList() {
  const navigate = useNavigate(); 

  const handleCreateContestClick = () => {
    navigate("/create-contest");
  }; 

  return (
    <div className="flex flex-col h-screen">
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
            {/* Rest of your cards */}
            <Card>
              <CardContent className="p-4 md:p-6 cursor-pointer">
                <div className="space-y-2">
                  <h3 className="text-lg font-semibold">CodeSprint</h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400">A 3-hour coding competition.</p>
                  <div className="flex items-center space-x-2 text-sm">
                    <ClockIcon className="w-4 h-4 text-gray-500" />
                    <time dateTime="2023-03-16T10:24:00">Starts at 10:24 AM</time>
                  </div>
                  <div className="flex items-center space-x-2 text-sm">
                    <ClockIcon className="w-4 h-4" />
                    <time dateTime="2023-03-16T13:24:00">Ends at 1:24 PM</time>
                  </div>
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 md:p-6 cursor-pointer">
                <div className="space-y-2">
                  <h3 className="text-lg font-semibold">CodeMaster</h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400">A 2-hour coding competition.</p>
                  <div className="flex items-center space-x-2 text-sm">
                    <ClockIcon className="w-4 h-4 text-gray-500" />
                    <time dateTime="2023-03-16T10:24:00">Starts at 10:24 AM</time>
                  </div>
                  <div className="flex items-center space-x-2 text-sm">
                    <ClockIcon className="w-4 h-4" />
                    <time dateTime="2023-03-16T13:24:00">Ends at 1:24 PM</time>
                  </div>
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 md:p-6 cursor-pointer">
                <div className="space-y-2">
                  <h3 className="text-lg font-semibold">HackDay</h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400">A 4-hour coding competition.</p>
                  <div className="flex items-center space-x-2 text-sm">
                    <ClockIcon className="w-4 h-4 text-gray-500" />
                    <time dateTime="2023-03-16T10:24:00">Starts at 10:24 AM</time>
                  </div>
                  <div className="flex items-center space-x-2 text-sm">
                    <ClockIcon className="w-4 h-4" />
                    <time dateTime="2023-03-16T13:24:00">Ends at 1:24 PM</time>
                  </div>
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 md:p-6 cursor-pointer">
                <div className="space-y-2">
                  <h3 className="text-lg font-semibold">AlgoQuest</h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400">A 1-hour coding competition.</p>
                  <div className="flex items-center space-x-2 text-sm">
                    <ClockIcon className="w-4 h-4 text-gray-500" />
                    <time dateTime="2023-03-16T10:24:00">Starts at 10:24 AM</time>
                  </div>
                  <div className="flex items-center space-x-2 text-sm">
                    <ClockIcon className="w-4 h-4" />
                    <time dateTime="2023-03-16T13:24:00">Ends at 1:24 PM</time>
                  </div>
                </div>
              </CardContent>
            </Card>
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
