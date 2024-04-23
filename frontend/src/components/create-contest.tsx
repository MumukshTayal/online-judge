import { useEffect, useState } from 'react';
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import Navbar from "./navbar.tsx"

export default function CreateContest() {
  const [formData, setFormData] = useState({
    name: '',
    startTime: '',
    endTime: '',
    description: '',
    problems: [],
    emails: [],
    languageLimits: [],
  });

  const [newProblem, setNewProblem] = useState('');
  const [newEmail, setNewEmail] = useState('');
  const [jwtToken, setJwtToken] = useState('');
  
  useEffect(() => {
    const storedToken = localStorage.getItem('jwtToken');
    if (storedToken) {
      setJwtToken(storedToken);
    }
  }, []);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleAddProblem = () => {
    if (newProblem.trim() === '') return; // prevent adding empty problem names
    setFormData({
      ...formData,
      problems: [...formData.problems, newProblem],
    });
    setNewProblem(''); // clear the input field after adding the problem
  };

  const handleRemoveProblem = (index) => {
    const updatedProblems = [...formData.problems];
    updatedProblems.splice(index, 1);
    setFormData({
      ...formData,
      problems: updatedProblems,
    });
  };

  const handleAddEmail = () => {
    if (newEmail.trim() === '') return; // prevent adding empty emails
    setFormData({
      ...formData,
      emails: [...formData.emails, newEmail],
    });
    setNewEmail(''); // clear the input field after adding the email
  };

  const handleRemoveEmail = (index) => {
    const updatedEmails = [...formData.emails];
    updatedEmails.splice(index, 1);
    setFormData({
      ...formData,
      emails: updatedEmails,
    });
  };

  const handleLanguageLimitChange = (index, type, value) => {
    const updatedLimits = [...formData.languageLimits];
    updatedLimits[index][type] = value;
    setFormData({
      ...formData,
      languageLimits: updatedLimits,
    });
  };

  const addLanguage = () => {
    const newLanguage = { id: Date.now(), language: '', timeLimit: '', memoryLimit: '' };
    setFormData({
      ...formData,
      languageLimits: [...formData.languageLimits, newLanguage],
    });
  };

  const removeLanguage = (id) => {
    const updatedLimits = formData.languageLimits.filter(limit => limit.id !== id);
    setFormData({
      ...formData,
      languageLimits: updatedLimits,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('http://localhost:8080/api/create_contest', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${jwtToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });
      // Handle response
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <Navbar />
        <div className="py-6 w-full space-y-6">
        <h1 className="text-3xl font-semibold tracking-tighter sm:text-4xl md:text-5xl">Create Contest</h1>
          <div>
            <div className="grid gap-4">
              <div className="grid gap-2">
                <Label htmlFor="name">Contest Name</Label>
                <Input id="name" name="name" placeholder="Enter contest name" required onChange={handleInputChange} />
              </div>
              <div className="grid gap-2 md:grid-cols-2">
                <div className="grid gap-2">
                  <Label htmlFor="start-time">Start Time</Label>
                  <Input id="start-time" name="startTime" placeholder="Enter start time" required type="datetime-local" onChange={handleInputChange} />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="end-time">End Time</Label>
                  <Input id="end-time" name="endTime" placeholder="Enter end time" required type="datetime-local" onChange={handleInputChange} />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="description">Contest Description</Label>
                <textarea id="description" name="description" className="w-full border rounded-lg px-3 py-2" rows="4" value={formData.description} onChange={handleInputChange} />
              </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <Label htmlFor="languages">Languages</Label>
              </div>
                {formData.languageLimits.map((limit, index) => (
                  <div key={limit.id} className="grid gap-2 md:grid-cols-4 items-center">
                    <div className="grid gap-2">
                      <Label htmlFor={`language-${index}`}>Language</Label>
                      <Select
                        id={`language-${index}`}
                        value={limit.language}
                        onValueChange={(value) => handleLanguageLimitChange(index, 'language', value)}
                      >
                        <SelectTrigger className="w-[180px]">
                          <SelectValue>{limit.language || "Select language"}</SelectValue>
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="java">Java</SelectItem>
                          <SelectItem value="python">Python</SelectItem>
                          <SelectItem value="c">C</SelectItem>
                          <SelectItem value="cpp">C++</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor={`time-limit-${index}`}>Time Limit (ms)</Label>
                      <Input id={`time-limit-${index}`} name={`timeLimit-${index}`} placeholder="Enter time limit" required type="number" value={limit.timeLimit} onChange={(e) => handleLanguageLimitChange(index, 'timeLimit', e.target.value)} />
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor={`memory-limit-${index}`}>Memory Limit (MB)</Label>
                      <Input id={`memory-limit-${index}`} name={`memoryLimit-${index}`} placeholder="Enter memory limit" required type="number" value={limit.memoryLimit} onChange={(e) => handleLanguageLimitChange(index, 'memoryLimit', e.target.value)} />
                    </div>
                    <div>
                      <Button onClick={() => removeLanguage(limit.id)}>Remove</Button>
                    </div>
                  </div>
                ))}
                <Button onClick={addLanguage}>Add Language</Button>
              </div>
              <div className="space-y-2">
                <Label htmlFor="problems">Problems</Label>
                <div className="flex flex-col w-full min-h-[200px] border rounded-lg">
                  <div className="grid w-full grid-cols-2 items-stretch divide-y p-2">
                    {formData.problems.map((problem, index) => (
                      <div key={index} className="flex w-full items-center justify-between px-2">
                        <Label className="text-sm font-medium leading-none" htmlFor={`problems-${index}`}>
                          {problem}
                        </Label>
                        <Button className="h-6 p-1 rounded-md" size="none" type="button" variant="ghost" onClick={() => handleRemoveProblem(index)}>
                          <XIcon className="w-4 h-4" />
                          <span className="sr-only">Remove</span>
                          <span className="text-lg leading-none -translate-y-px-5">
                            <XIcon className="w-4 h-4 inline-block" />
                          </span>
                        </Button>
                      </div>
                    ))}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Input className="max-w-sm flex-1" placeholder="Enter problem name" type="text" value={newProblem} onChange={(e) => setNewProblem(e.target.value)} />
                  <Button type="button" onClick={handleAddProblem}>+</Button>
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="emails">Emails</Label>
                <div className="flex space-x-2">
                  <Input className="max-w-sm flex-1" id="emails" placeholder="Enter email" type="email" value={newEmail} onChange={(e) => setNewEmail(e.target.value)} />
                  <Button type="button" onClick={handleAddEmail}>Add</Button>
                </div>
                <div>
                  {formData.emails.map((email, index) => (
                    <div key={index} className="flex items-center space-x-2">
                      <UserIcon className="w-4 h-4" />
                      <span className="text-sm font-medium">{email}</span>
                      <Button size="none" variant="ghost" onClick={() => handleRemoveEmail(index)}>
                        <XIcon className="w-4 h-4" />
                        <span className="sr-only">Remove</span>
                        <XIcon className="w-4 h-4 inline-block" />
                      </Button>
                    </div>
                  ))}
                </div>
              </div>
              <div className="flex w-full pt-4">
                <Button className="ml-auto" type="submit">Create Contest</Button>
              </div>
            </div>
          </div>
        </div> 
      </div>
    </form>
  );
}

function UserIcon(props) {
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
      <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  );
}

function XIcon(props) {
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
      <path d="M18 6 6 18" />
      <path d="M6 6 18 18" />
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
