import { useState, useEffect } from 'react';
import { CardTitle, CardDescription, CardHeader, CardFooter, Card } from "@/components/ui/card";
import Navbar from "./navbar.tsx"

export default function HomePage() {
  const [contests, setContests] = useState([]);

  useEffect(() => {
    // Fetch data from the API
    fetch('http://localhost:8080/api/get_all_contests')
      .then(response => response.json())
      .then(data => {
        // Set the fetched contests to state
        setContests(data);
      })
      .catch(error => console.error('Error fetching contests:', error));
  }, []);

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center space-y-4 text-center">
              <div className="space-y-2">
                <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl lg:text-6xl/none">
                  Welcome to Online Judge
                </h1>
                <p className="mx-auto max-w-[700px] text-gray-500 md:text-xl dark:text-gray-400">
                  Join our community and improve your coding skills.
                </p>
              </div>
              {/* <div className="space-x-4">
                <a
                  className="inline-flex h-9 items-center justify-center rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-gray-50 shadow transition-colors hover:bg-gray-900/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-gray-950 disabled:pointer-events-none disabled:opacity-50 dark:bg-gray-50 dark:text-gray-900 dark:hover:bg-gray-50/90 dark:focus-visible:ring-gray-300"
                  href="#"
                >
                  Sign Up
                </a>
                <a
                  className="inline-flex h-9 items-center justify-center rounded-md border border-gray-200 border-gray-200 bg-white px-4 py-2 text-sm font-medium shadow-sm transition-colors hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-gray-950 disabled:pointer-events-none disabled:opacity-50 dark:border-gray-800 dark:border-gray-800 dark:bg-gray-950 dark:hover:bg-gray-800 dark:hover:text-gray-50 dark:focus-visible:ring-gray-300"
                  href="#"
                >
                  Log In
                </a>
              </div> */}
            </div>
          </div>
        </section>
        <section className="w-full py-12 md:py-24 lg:py-32 bg-gray-100 dark:bg-gray-800">
          <div className="container px-4 md:px-6">
            <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">Upcoming Contests</h2>
            <div className="grid gap-6 mt-6 md:grid-cols-2 lg:grid-cols-3">
              {/* Map through contests and display each one */}
              {contests.map(contest => (
                <Card key={contest.contest_id}>
                  <CardHeader>
                    <CardTitle>{contest.contest_title}</CardTitle>
                    <CardDescription>Starts: {new Date(contest.start_time).toLocaleString()}</CardDescription>
                  </CardHeader>
                  <CardFooter>
                    <a
                      className="inline-flex h-9 items-center justify-center rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-gray-50 shadow transition-colors hover:bg-gray-900/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-gray-950 disabled:pointer-events-none disabled:opacity-50 dark:bg-gray-50 dark:text-gray-900 dark:hover:bg-gray-50/90 dark:focus-visible:ring-gray-300"
                      href={`/contest/${contest.contest_id}`}
                    >
                      View Details
                    </a>
                  </CardFooter>
                </Card>
              ))}
            </div>
          </div>
        </section>
        <section className="w-full py-12 md:py-24 lg:py-32">
          <footer className="py-6 bg-gray-900 text-gray-50 text-center">© 2024 Online Judge.</footer>
        </section>
      </main>
    </div>
  );
}
