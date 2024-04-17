import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Avatar } from "@/components/ui/avatar";
import { Textarea } from "@/components/ui/textarea";
import { SelectValue, SelectTrigger, SelectItem, SelectContent, Select } from "@/components/ui/select";
import Navbar from "./navbar.tsx"

export default function ProblemView() {
  const { problemId } = useParams();
  const [problem, setProblem] = useState(null);

  useEffect(() => {
    const fetchProblem = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/get_problem/${problemId}`);
        if (!response.ok) {
          throw new Error('Failed to fetch problem');
        }
        const data = await response.json();
        setProblem(data);
      } catch (error) {
        console.error(error);
      }
    };

    fetchProblem();
  }, [problemId]);

  if (!problem) {
    return <div>Loading...</div>;
  }

  return (
    <div className="w-full px-4 py-6 space-y-6 md:px-6">
      <Navbar />
      <div className="space-y-2">
        <h1 className="text-3xl font-bold">{problem.ProblemTitle}</h1>
        <p className="text-gray-500 dark:text-gray-400">
          {problem.ProblemDescription}
        </p>
      </div>
      <div className="space-y-2">
        <h2 className="text-2xl font-bold">Constraints</h2>
        <p className="text-gray-500 dark:text-gray-400">
          {problem.ConstraintsDesc}
        </p>
      </div>
      <div className="space-y-2">
        <h2 className="text-2xl font-bold">Input</h2>
        <p className="text-gray-500 dark:text-gray-400">{problem.InputFormat}</p>
      </div>
      <div className="space-y-2">
        <h2 className="text-2xl font-bold">Output</h2>
        <p className="text-gray-500 dark:text-gray-400">{problem.OutputFormat}</p>
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        <div>
          <h2 className="text-2xl font-bold">Sample Input</h2>
          <p>{problem.SampleInput}</p>
        </div>
        <div>
          <h2 className="text-2xl font-bold">Sample Output</h2>
          <p>{problem.SampleOutput}</p>
        </div>
      </div>
      <div className="flex items-center space-x-4">
        <div className="flex items-center space-x-2">
          <div>
            <h3 className="font-semibold">{problem.CreatorName}</h3>
            <p className="text-sm text-gray-500 dark:text-gray-400">Author: {problem.CreatorEmail}</p>
          </div>
        </div>
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        <div>
          <h2 className="text-2xl font-bold">Add Your Code</h2>
          <Textarea placeholder="Enter your code here." />
        </div>
        <div>
          <h2 className="text-2xl font-bold">Select Language</h2>
          <Select>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Select language" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="javascript">C</SelectItem>
              <SelectItem value="python">Python</SelectItem>
              <SelectItem value="java">Java</SelectItem>
              <SelectItem value="typescript">C++</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>
      <div className="flex justify-start space-x-4">
        <button className="bg-white hover:bg-gray-200 text-black font-bold py-2 px-4 rounded">
          Run
        </button>
        <button className="bg-black hover:bg-gray-800 text-white font-bold py-2 px-4 rounded">
          Submit
        </button>
      </div>
    </div>
  )
}
