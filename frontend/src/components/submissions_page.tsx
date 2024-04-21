import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { CardTitle, CardHeader, CardContent, Card } from "@/components/ui/card"
import { TableHead, TableRow, TableHeader, TableCell, TableBody, Table } from "@/components/ui/table"
import Navbar from "./navbar.tsx"

export default function SubmissionsPage() {
  const { contestId } = useParams();
  const [submissions, setSubmissions] = useState([]);

  useEffect(() => {
    // Fetch submissions data from the backend API
    fetch(`http://localhost:8080/api/get_submissions?contestId=${contestId}`)
      .then(response => response.json())
      .then(data => {
        // Convert object to array of submissions
        const formattedSubmissions = Object.values(data).flatMap(submissionArray => submissionArray);
        setSubmissions(formattedSubmissions);
      })
      .catch(error => console.error('Error fetching submissions:', error));
  }, []);

  return (
    <div>
      <Navbar />
      <b className="text-3xl font-semibold tracking-tighter sm:text-4xl md:text-5xl">Submissions for Contest {contestId}</b>
      <div className="flex justify-center min-h-screen bg-gray-100 dark:bg-gray-800 mt-4">
        <div className="max-w-3xl w-full p-4">
          <Card className="h-full">
            <CardHeader>
              <CardTitle className="text-lg">Submissions</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              <div className="border-t border-gray-200 dark:border-gray-800">
                <div className="relative w-full overflow-auto">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        {/* <TableHead className="w-[120px]">Submission ID</TableHead> */}
                        <TableHead>User</TableHead>
                        <TableHead>Verdict</TableHead>
                        <TableHead>Time</TableHead>
                        <TableHead>Language</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {/* Map through submissions and display each submission */}
                      {submissions.map((submission, index) => (
                        <TableRow key={index}>
                          <TableCell>{submission.UserID}</TableCell>
                          <TableCell>{submission.Result}</TableCell>
                          <TableCell>{submission.Time}</TableCell>
                          <TableCell>{submission.Language}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
