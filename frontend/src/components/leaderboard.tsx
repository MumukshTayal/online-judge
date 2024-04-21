import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Navbar from "./navbar.tsx"

export default function LeaderboardPage() {
  const { contestId } = useParams();
  const [leaderboard, setLeaderboard] = useState([]);

  useEffect(() => {
    fetch(`http://localhost:8080/api/get_leaderboard?contestId=${contestId}`)
      .then(response => response.json())
      .then(data => {
        setLeaderboard(data);
      })
      .catch(error => console.error('Error fetching leaderboard:', error));
  }, [contestId]);

    return (
        <div>
            <Navbar />
    <div className="flex flex-col min-h-screen">
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32">
              <div className="container px-4 md:px-6">
              <b className="text-3xl font-semibold tracking-tighter sm:text-4xl md:text-5xl">Leaderboard for Contest {contestId}</b>
            <div className="overflow-x-auto mt-6">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      User Email
                    </th>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Score
                    </th>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Last Submission Time
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {/* Map through leaderboard and display each user */}
                  {leaderboard.map(entry => (
                    <tr key={entry.UserEmail}>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {entry.UserEmail}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {entry.Result}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {entry.LastSubmissionTime}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </section>
      </main>
            </div>
            </div>
  );
}
