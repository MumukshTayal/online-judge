/**
 * This code was generated by v0 by Vercel.
 * @see https://v0.dev/t/HLXBT8sbuRQ
 */
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"

export function AddTestCasePage() {
  return (
    <div className="mx-auto max-w-3xl px-4">
      <div className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="problem">Problem</Label>
          <Input id="problem" placeholder="Enter the problem name" type="text" />
        </div>
        <div className="space-y-2">
          <Label>Input</Label>
          <Input accept=".txt" id="input" type="file" />
          <div>Upload the input file for the test case. It should be named "input.txt".</div>
        </div>
        <div className="space-y-2">
          <Label>Output</Label>
          <Input accept=".txt" id="output" type="file" />
          <div>Upload the output file for the test case. It should be named "output.txt".</div>
        </div>
        <Button>Upload</Button>
      </div>
    </div>
  )
}
