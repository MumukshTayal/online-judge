import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Navbar from "./navbar.tsx"

export default function AllSubmissionsPage() {
  const [submissions, setSubmissions] = useState([]);

  useEffect(() => {
    // Fetch submissions data from the backend API
    fetch(`http://localhost:8080/api/get_all_submissions`)
      .then(response => response.json())
      .then(data => {
        // Convert object to array of submissions (replace with your logic)
        const formattedSubmissions = Object.values(data).flatMap(submissionArray => submissionArray);
        setSubmissions(formattedSubmissions);
      })
      .catch(error => console.error('Error fetching submissions:', error));
  }, []);

  return (
      <div>
          <Navbar/>
      <table style={{ width: '100%', borderCollapse: 'collapse', padding: '32 px'}}>
        <thead>
          <tr style={{ backgroundColor: '#f5f5f5' }}>
          <th style={{ padding: '8px', border: '1px solid #ddd' }}>User</th>
            <th style={{ padding: '8px', border: '1px solid #ddd' }}>Verdict</th>
            <th style={{ padding: '8px', border: '1px solid #ddd' }}>Time</th>
            <th style={{ padding: '8px', border: '1px solid #ddd' }}>Language</th>
          </tr>
        </thead>
        <tbody>
          {submissions.map((submission, index) => (
            <tr key={index}>
            <td style={{ padding: '8px', border: '1px solid #ddd' }}>{submission.UserID}</td>
              <td style={{ padding: '8px', border: '1px solid #ddd' }}>{submission.Result}</td>
              <td style={{ padding: '8px', border: '1px solid #ddd' }}>{submission.Time}</td>
              <td style={{ padding: '8px', border: '1px solid #ddd' }}>{submission.Language}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
