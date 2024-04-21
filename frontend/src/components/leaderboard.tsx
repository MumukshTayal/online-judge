import { useState, useEffect } from 'react';
import { useNavigate, Route, useParams } from 'react-router-dom';
import { CardTitle, CardDescription, CardHeader, CardFooter, Card } from "@/components/ui/card";
import Navbar from "./navbar.tsx"

export default function LeaderboardPage() {
  const { contestId } = useParams();
  const [leaderboard, setLeaderboard] = useState([]);

  useEffect(() => {
    // Fetch data from the API
    fetch('http://localhost:8080/api/get_leaderboard')
      .then(response => response.json())
      .then(data => {
        // Set the fetched leaderboard to state
        setLeaderboard(data);
      })
      .catch(error => console.error('Error fetching leaderboard:', error));
  }, []);

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32">
          <div className="container px-4 md:px-6">
            <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl lg:text-6xl/none">
              Leaderboard for Contest {contestId}
            </h1>
            <div className="grid gap-6 mt-6 md:grid-cols-2 lg:grid-cols-3">
              {/* Map through leaderboard and display each user */}
              {leaderboard.map(user => (
                <Card key={user.user_id}>
                  <CardHeader>
                    <CardTitle>{user.username}</CardTitle>
                    <CardDescription>Score: {user.score}</CardDescription>
                  </CardHeader>
                  <CardFooter>
                    <p className="text-gray-500">Rank: {user.rank}</p>
                  </CardFooter>
                </Card>
              ))}
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
