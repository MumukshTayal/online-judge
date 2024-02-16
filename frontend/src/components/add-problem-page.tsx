/**
 * This code was generated by v0 by Vercel.
 * @see https://v0.dev/t/l1iiotDPjgg
 */
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { SelectValue, SelectTrigger, SelectItem, SelectContent, Select } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { Button } from "@/components/ui/button"
import { CardTitle, CardDescription, CardHeader, CardContent, Card } from "@/components/ui/card"

export function AddProblemPage() {
  return (
    <div key="1" className="grid max-w-6xl w-full gap-6 p-4 mx-auto lg:grid-cols-3 lg:gap-10">
      <div className="space-y-4 lg:col-span-2">
        <div className="space-y-2">
          <h1 className="text-3xl font-bold">Create a new problem for your contest</h1>
          <p className="text-gray-500 grid-cols-2 gap-2 dark:text-gray-400">
            Add a new problem to your contest so that your participants can submit their solutions.
          </p>
        </div>
        <div className="space-y-2">
          <Label htmlFor="title">Title</Label>
          <Input id="title" placeholder="Sum of two numbers" required />
        </div>
        <div className="space-y-2">
          <Label htmlFor="slug">Slug</Label>
          <Input id="slug" placeholder="sum-of-two-numbers" required />
          <p className="text-sm text-gray-500 dark:text-gray-400">
            The slug is the unique URL-friendly identifier for your problem.
          </p>
        </div>
        <div className="space-y-2">
          <Label htmlFor="difficulty">Difficulty</Label>
          <Select defaultValue="easy">
            <SelectTrigger className="w-36">
              <SelectValue placeholder="Select" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="easy">Easy</SelectItem>
              <SelectItem value="medium">Medium</SelectItem>
              <SelectItem value="hard">Hard</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="statement">Problem Statement</Label>
          <Textarea
            className="min-h-[200px]"
            id="statement"
            placeholder="Write your problem statement in Markdown or LaTeX"
            required
          />
        </div>
        <div className="space-y-2">
          <Label>Input Format</Label>
          <Textarea
            className="min-h-[200px]"
            id="input"
            placeholder="Write the input format in Markdown or LaTeX"
            required
          />
        </div>
        <div className="space-y-2">
          <Label>Output Format</Label>
          <Textarea
            className="min-h-[200px]"
            id="output"
            placeholder="Write the output format in Markdown or LaTeX"
            required
          />
        </div>
        <div className="space-y-2">
          <Label>Constraints</Label>
          <Textarea
            className="min-h-[200px]"
            id="constraints"
            placeholder="Write the constraints in Markdown or LaTeX"
            required
          />
        </div>
        <div className="space-y-2">
          <Label>Sample Input</Label>
          <Textarea
            className="min-h-[200px]"
            id="sample-input"
            placeholder="Write the sample input in Markdown or LaTeX"
            required
          />
        </div>
        <div className="space-y-2">
          <Label>Sample Output</Label>
          <Textarea
            className="min-h-[200px]"
            id="sample-output"
            placeholder="Write the sample output in Markdown or LaTeX"
            required
          />
        </div>
        <div className="space-y-2">
          <Label>Test Cases</Label>
          <Textarea
            className="min-h-[200px]"
            id="test-cases"
            placeholder="Add your test cases in the format input1,input2|output1,output2"
            required
          />
        </div>
        <div className="flex flex-col gap-2 min-[400px]:flex-row justify-center">
          <Button className="w-full" type="submit">
            Create Problem
          </Button>
        </div>
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
            <li>Add test cases in the 'Test Cases' textarea.</li>
            <li>Click the 'Create Problem' button to finalize the process.</li>
          </ol>
        </CardContent>
      </Card>
    </div>
  )
}