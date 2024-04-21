import React, { useEffect, useState, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { Textarea } from "@/components/ui/textarea";
import Navbar from "./navbar.tsx"

export default function ProblemView() {
  const { contestId } = useParams();
  const { problemId } = useParams();
  const [problem, setProblem] = useState(null);
  const [code, setCode] = useState("");
  const [language, setLanguage] = useState('py');
  const [jwtToken, setJwtToken] = useState('');
  const codeTextareaRef = useRef(null);
  const [result, setResult] = useState("");

  useEffect(() => {
    const storedToken = localStorage.getItem('jwtToken');
    if (storedToken) {
      setJwtToken(storedToken);
    }
  }, []);

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

  const handleRunClick = () => {
    // if (!language) {
    //   console.error('Please select a language before running.');
    //   return;
    // }
    const requestBody = {
      contest_id: contestId,
      problem_id: problemId,
      code: codeTextareaRef.current.value,
      language: language
    };

    fetch('http://localhost:8080/api/run_code', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${jwtToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    })
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to run code');
      }
      return response.text()
    })
      .then(result => {
        console.log("Result:", result)
        setResult(result);
    })
    .catch(error => {
      console.error('Error running code:', error);
    });
  };

  const handleSubmit = () => {
    // if (!language) {
    //   console.error('Please select a language before running.');
    //   return;

    const requestBody = {
      contest_id: contestId,
      problem_id: problemId,
      code: codeTextareaRef.current.value,
      language: language,
      submission_date_time: new Date().toISOString()
    };

    fetch('http://localhost:8080/api/submit', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${jwtToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    })
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to submit code');
      }
      return response.text()
    })
      .then(result => {
        console.log("Result:", result)
        setResult(result);
    })
    .catch(error => {
      console.error('Error submitting code:', error);
    });
  };

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
          {/* <p>{problem.SampleInput}</p> */}
          <div dangerouslySetInnerHTML={{ __html: problem.SampleInput.replace(/\n/g, '<br>') }} />
        </div>
        <div>
          <h2 className="text-2xl font-bold">Sample Output</h2>
          {/* <p>{problem.SampleOutput}</p> */}
          <div dangerouslySetInnerHTML={{ __html: problem.SampleOutput.replace(/\n/g, '<br>') }} />
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
          <Textarea ref={codeTextareaRef} placeholder="Enter your code here." />
        </div>
        <div>
          <h2 className="text-2xl font-bold">Select Language</h2>
          <select onChange={(e) => setLanguage(e.target.value)}>
            <option value="py">Python</option>
            <option value="cpp">C++</option>
            <option value="c">C</option>
          </select>
        </div>
      </div>
      <div className="flex justify-start space-x-4">
        <button onClick={handleRunClick} className="bg-white hover:bg-gray-200 text-black font-bold py-2 px-4 rounded">
          Run
        </button>
        <button onClick={handleSubmit} className="bg-black hover:bg-gray-800 text-white font-bold py-2 px-4 rounded">
          Submit
        </button>
      </div>
      <div>
          <h2 className="text-2xl font-bold">Result</h2>
        {/* <p className="text-gray-500 dark:text-gray-400">{result}</p> */}
        <div dangerouslySetInnerHTML={{ __html: result.replace(/\n/g, '<br>') }} />
        </div>
    </div>
  )
}
