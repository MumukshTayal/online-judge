import React from 'react';

function Navbar() {
  return (
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
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/problems">
            Problems
          </a>
          <a className="text-sm font-medium hover:underline underline-offset-4" href="/submissions">
            Submissions 
          </a>
          {/* <a className="text-sm font-medium hover:underline underline-offset-4" href="/add-testcase">
            Add Test Cases 
          </a> */}
        </nav>
      </header>
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
  );
}

export default Navbar;
