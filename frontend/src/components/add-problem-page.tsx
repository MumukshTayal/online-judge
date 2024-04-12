import { useState, useEffect } from 'react';
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { SelectValue, SelectTrigger, SelectItem, SelectContent, Select } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { Button } from "@/components/ui/button"
import { CardTitle, CardDescription, CardHeader, CardContent, Card } from "@/components/ui/card"
import Navbar from "./navbar.tsx"

export default function AddProblemPage() { 
  const [jwtToken, setJwtToken] = useState(''); 
  useEffect(() => {
    const storedToken = localStorage.getItem('jwtToken');
    if (storedToken) {
      setJwtToken(storedToken);
    }
  }, []);

  const [formData, setFormData] = useState({
    title: '',
    statement: '',
    constraints: '',
    sampleInput: '',
    input: '',
    output: '',
    sampleOutput: '',
    isPrivate: '' 
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log("type of isPrivate --> ", typeof(formData.isPrivate));
    try { 
      const response = await fetch('http://localhost:8080/api/create_problem', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${jwtToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });
      if (response.ok) {
        // Problem created successfully, perform any necessary actions
        console.log('Problem created successfully');
      } else {
        // Handle error
        console.error('Problem creation failed');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <div>
      <Navbar />
      <div key="1" className="grid max-w-6xl w-full gap-6 p-4 mx-auto lg:grid-cols-3 lg:gap-10">
        <div className="space-y-4 lg:col-span-2">
          <div className="space-y-2">
            <h1 className="text-3xl font-bold">Create a new problem for your contest</h1>
            <p className="text-gray-500 grid-cols-2 gap-2 dark:text-gray-400">
              Add a new problem to your contest so that your participants can submit their solutions.
            </p>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label htmlFor="title">Title</Label>
              <Input id="title" name="title" placeholder="Sum of two numbers" required onChange={handleInputChange} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="statement">Problem Statement</Label>
              <Textarea
                className="min-h-[200px]"
                id="statement"
                name="statement"
                placeholder="Write your problem statement in Markdown or LaTeX"
                required
                onChange={handleInputChange}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="input">Input Format</Label>
              <Textarea
                className="min-h-[200px]"
                id="input"
                name="input"
                placeholder="Write the input format in Markdown or LaTeX"
                required
                onChange={handleInputChange}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="output">Output Format</Label>
              <Textarea
                className="min-h-[200px]"
                id="output"
                name="output"
                placeholder="Write the output format in Markdown or LaTeX"
                required
                onChange={handleInputChange}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="constraints">Constraints</Label>
              <Textarea
                className="min-h-[200px]"
                id="constraints"
                name="constraints"
                placeholder="Write the constraints in Markdown or LaTeX"
                required
                onChange={handleInputChange}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="sample-input">Sample Input</Label>
              <Textarea
                className="min-h-[200px]"
                id="sample-input"
                name="sampleInput"
                placeholder="Write the sample input in Markdown or LaTeX"
                required
                onChange={handleInputChange}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="sample-output">Sample Output</Label>
              <Textarea
                className="min-h-[200px]"
                id="sample-output"
                name="sampleOutput"
                placeholder="Write the sample output in Markdown or LaTeX"
                required
                onChange={handleInputChange}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="is-private">Is Private</Label> 
              <select id="is-private" name="isPrivate" onChange={handleInputChange}>
                <option value='true'>Yes</option>
                <option value='false'>No</option>
              </select>
            </div>
            <div className="flex flex-col gap-2 min-[400px]:flex-row justify-center">
              <Button className="w-full" type="submit">
                Create Problem
              </Button>
            </div>
          </form>
        </div>
        <Card className="p-4 space-y-4">
          <CardHeader className="pb-0">
            <CardTitle>How to Add a Problem</CardTitle>
            <CardDescription>Learn how to add a new problem to your contest.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <p>To add a new problem to your contest, follow these steps:</p>
            <ol className="list-decimal pl-4">
              <li>Enter the title of the problem in the 'Title' field.</li>
              <li>Provide a unique slug for the problem in the 'Slug' field.</li>
              <li>Select the difficulty level from the dropdown menu.</li>
              <li>Write the problem statement in the 'Problem Statement' textarea.</li>
              <li>Describe the input format in the 'Input Format' textarea.</li>
              <li>Define the output format in the 'Output Format' textarea.</li>
              <li>Add any constraints in the 'Constraints' textarea.</li>
              <li>Include sample input in the 'Sample Input' textarea.</li>
              <li>Specify the expected output in the 'Sample Output' textarea.</li>
              <li>Enter your email in the 'Creator Email' field.</li>
              <li>Choose whether the problem is private or not.</li>
              <li>Click the 'Create Problem' button to finalize the process.</li>
              <li>Add test cases using the 'Add Test Cases' page.</li>
            </ol>
          </CardContent>
        </Card>
      </div>
    </div>
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
