import { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

export default function CreateContest() {
  const [problems, setProblems] = useState([]);
  const [newProblem, setNewProblem] = useState('');
  const [emails, setEmails] = useState([]);
  const [newEmail, setNewEmail] = useState('');
  const [contestDescription, setContestDescription] = useState('');
  const [languageLimits, setLanguageLimits] = useState([]);

  const handleAddProblem = () => {
    if (newProblem.trim() === '') return; // prevent adding empty problem names
    setProblems([...problems, newProblem]);
    setNewProblem(''); // clear the input field after adding the problem
  };

  const handleRemoveProblem = (index) => {
    const updatedProblems = [...problems];
    updatedProblems.splice(index, 1);
    setProblems(updatedProblems);
  };

  const handleAddEmail = () => {
    if (newEmail.trim() === '') return; // prevent adding empty emails
    setEmails([...emails, newEmail]);
    setNewEmail(''); // clear the input field after adding the email
  };

  const handleRemoveEmail = (index) => {
    const updatedEmails = [...emails];
    updatedEmails.splice(index, 1);
    setEmails(updatedEmails);
  };

  const handleLanguageLimitChange = (index, type, value) => {
    const updatedLimits = [...languageLimits];
    updatedLimits[index][type] = value;
    setLanguageLimits(updatedLimits);
  };

  const addLanguage = () => {
    const newLanguage = { id: Date.now(), language: '', timeLimit: '', memoryLimit: '' };
    setLanguageLimits([...languageLimits, newLanguage]);
  };

  const removeLanguage = (id) => {
    const updatedLimits = languageLimits.filter(limit => limit.id !== id);
    setLanguageLimits(updatedLimits);
  }; 


  return (
    <div>
      <header className="px-4 lg:px-6 h-14 flex items-center">
        <nav>
          <a className="flex items-center justify-center" href="/home">
            <CodeIcon className="h-6 w-6" />
            <span className="sr-only">Online Judge</span>
          </a> 
        </nav>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/contest-list">
            Contests
          </a>
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/add-problem">
            Add Problem 
          </a>
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/submissions">
            Submissions 
          </a>
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/add-testcase">
            Add Test Cases 
          </a>
        </nav>
      </header> 
      <div className="py-6 w-full space-y-6">
        <div>
          <div className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Contest Name</Label>
              <Input id="name" placeholder="Enter contest name" required />
            </div>
            <div className="grid gap-2 md:grid-cols-2">
              <div className="grid gap-2">
                <Label htmlFor="start-time">Start Time</Label>
                <Input id="start-time" placeholder="Enter start time" required type="datetime-local" />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="end-time">End Time</Label>
                <Input id="end-time" placeholder="Enter end time" required type="datetime-local" />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Contest Description</Label>
              <textarea id="description" className="w-full border rounded-lg px-3 py-2" rows="4" value={contestDescription} onChange={(e) => setContestDescription(e.target.value)} />
            </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="languages">Languages</Label>
            </div>
              {languageLimits.map((limit, index) => (
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
                    <Input id={`time-limit-${index}`} placeholder="Enter time limit" required type="number" value={limit.timeLimit} onChange={(e) => handleLanguageLimitChange(index, 'timeLimit', e.target.value)} />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor={`memory-limit-${index}`}>Memory Limit (MB)</Label>
                    <Input id={`memory-limit-${index}`} placeholder="Enter memory limit" required type="number" value={limit.memoryLimit} onChange={(e) => handleLanguageLimitChange(index, 'memoryLimit', e.target.value)} />
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
                  {problems.map((problem, index) => (
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
                {emails.map((email, index) => (
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
              <Button className="ml-auto">Create Contest</Button>
            </div>
          </div>
        </div>
      </div> 
    </div>
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
