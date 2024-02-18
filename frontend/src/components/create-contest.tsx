/**
 * This code was generated by v0 by Vercel.
 * @see https://v0.dev/t/LnFB3zQJO6N
 */
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"

export function CreateContest() {
  return (
    <div className="py-6 w-full space-y-6">
      <div>
        <div className="flex items-center space-x-4">
          <h1 className="text-3xl font-bold tracking-tighter">Create Contest</h1>
          <Button size="sm" variant="outline">
            <UserIcon className="w-4 h-4" />
            Invite Teammates{"\n                  "}
          </Button>
        </div>
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
            <Label htmlFor="problems">Problems</Label>
            <div className="flex flex-col w-full min-h-[200px] border rounded-lg">
              <div className="grid w-full grid-cols-2 items-stretch divide-y p-2">
                <div className="flex w-full items-center justify-between px-2">
                  <Label className="text-sm font-medium leading-none" htmlFor="problems-0">
                    A. Two Sum
                  </Label>
                  <Button className="h-6 p-1 rounded-md" size="none" type="button" variant="ghost">
                    <XIcon className="w-4 h-4" />
                    <span className="sr-only">Remove</span>
                    <span className="text-lg leading-none -translate-y-px-5">
                      <XIcon className="w-4 h-4 inline-block" />
                    </span>
                  </Button>
                </div>
                <div className="flex w-full items-center justify-between px-2">
                  <Label className="text-sm font-medium leading-none" htmlFor="problems-1">
                    B. Add Two Numbers
                  </Label>
                  <Button className="h-6 p-1 rounded-md" size="none" type="button" variant="ghost">
                    <XIcon className="w-4 h-4" />
                    <span className="sr-only">Remove</span>
                    <span className="text-lg leading-none -translate-y-px-5">
                      <XIcon className="w-4 h-4 inline-block" />
                    </span>
                  </Button>
                </div>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Input className="max-w-sm flex-1" placeholder="Enter problem name" type="text" />
              <Button type="button">+</Button>
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="emails">Emails</Label>
            <div className="flex space-x-2">
              <Input className="max-w-sm flex-1" id="emails" placeholder="Enter email" type="email" />
              <Button type="button">Add</Button>
            </div>
            <div className="flex items-center space-x-2">
              <UserIcon className="w-4 h-4" />
              <span className="text-sm font-medium">alice@example.com</span>
              <Button size="none" variant="ghost">
                <XIcon className="w-4 h-4" />
                <span className="sr-only">Remove</span>
                <XIcon className="w-4 h-4 inline-block" />
              </Button>
            </div>
          </div>
          <div className="flex w-full pt-4">
            <Button className="ml-auto">Create Contest</Button>
          </div>
        </div>
      </div>
    </div>
  )
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
  )
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
      <path d="m6 6 12 12" />
    </svg>
  )
}