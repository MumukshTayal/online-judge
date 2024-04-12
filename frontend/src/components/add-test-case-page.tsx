import { useState, useEffect } from 'react';
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import Navbar from "./navbar.tsx";

export default function AddTestCasePage() {
  const [problemName, setProblemName] = useState('');
  const [inputFile, setInputFile] = useState(null);
  const [outputFile, setOutputFile] = useState(null); 
  const [jwtToken, setJwtToken] = useState(''); 

  useEffect(() => {
    const storedToken = localStorage.getItem('jwtToken');
    if (storedToken) {
      setJwtToken(storedToken);
    }
  }, []);

  const handleInputChange = (event) => {
    const file = event.target.files[0];
    setInputFile(file);
  };

  const handleOutputChange = (event) => {
    const file = event.target.files[0];
    setOutputFile(file);
  };

  const handleUpload = async () => {
    try {
      const formData = new FormData();
      formData.append('problem_name', problemName);
      formData.append('input', inputFile);
      formData.append('output', outputFile);

      const response = await fetch('http://localhost:8080/api/add_testcase', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${jwtToken}`,
        },
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Failed to add test case');
      }

      alert('Test case added successfully!');
    } catch (error) {
      console.error(error);
      alert('Failed to add test case');
    }
  };

  return ( 
    <div>
      <Navbar />
      <div className="mx-auto max-w-3xl px-4">
        <div className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="problem">Problem Name</Label>
            <Input id="problem" placeholder="Enter the problem name" type="text" value={problemName} onChange={(e) => setProblemName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Input</Label>
            <Input accept=".txt" id="input" type="file" onChange={handleInputChange} />
            <div>Upload the input file for the test case. It should be named "input.txt".</div>
          </div>
          <div className="space-y-2">
            <Label>Output</Label>
            <Input accept=".txt" id="output" type="file" onChange={handleOutputChange} />
            <div>Upload the output file for the test case. It should be named "output.txt".</div>
          </div>
          <Button onClick={handleUpload}>Upload</Button>
        </div>
      </div> 
    </div>
  );
}
