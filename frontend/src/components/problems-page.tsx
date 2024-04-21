import { useState, useEffect } from 'react';
import { CardContent, Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from 'react-router-dom';
import Navbar from "./navbar.tsx";

export default function ProblemsPage() {
  const [problems, setProblems] = useState([]);

  useEffect(() => {
    fetchProblems();
  }, []);

  const fetchProblems = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/get_all_problems");
      if (!response.ok) {
        throw new Error('Failed to fetch problems');
      }
      const data = await response.json();
      setProblems(data);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="flex flex-col h-screen">
      <Navbar />
      <main className="flex-1 overflow-y-auto dark:bg-gray-900">
        <section className="container py-6 space-y-6 text-gray-900 md:space-y-8 dark:text-gray-50">
          <div className="space-y-6">
            <b className="text-3xl font-semibold tracking-tighter sm:text-4xl md:text-5xl">Contribute</b>
            <p className="max-w text-gray-500 md:text-base/relaxed dark:text-gray-400">
              Add a new problem to our database to challenge our participants and foster creativity in problem-solving. Your contribution will enrich our platform, providing diverse and engaging challenges for our community of learners and enthusiasts.
            </p>
            <Link to="/add-problem">
              <Button variant="primary" className="w-32 bg-black text-white mt-4">
                Add Problem
              </Button>
            </Link>
            <p className="max-w text-gray-500 md:text-base/relaxed dark:text-gray-400">
            Add new testcases to enrich our database and improve problem evaluation.
            </p>
            <Link to="/add-testcase">
              <Button variant="primary" className="w-32 bg-black text-white mt-4">
                Add Testcase
              </Button>
            </Link>
          </div>

          <div className="mt-16">
            <h1 className="text-3xl font-semibold tracking-tighter sm:text-4xl md:text-5xl">All Problems</h1>
            <div className="overflow-x-auto mt-4">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50 dark:bg-gray-800">
                  <tr>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Problem ID
                    </th>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Problem Title
                    </th>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Author
                    </th>
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Public
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {problems.map((problem, index) => (
                    <tr key={index}>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900">{problem.ProblemID}</div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900">{problem.ProblemTitle}</div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900">{problem.CreatorEmail}</div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900">{problem.IsPrivate ? 'No' : 'Yes'}</div>
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
  );
}
