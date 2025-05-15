import { useEffect, useState } from 'react';
import axios from 'axios';
import JobChart from './components/JobChart';

export interface Job {
  id: string;
  repo_url: string;
  branch: string;
  status: string;
  created_at: string;
  finished_at: string;
  lint_errors: number;
  lint_warnings: number;
}

function App() {
  const [jobs, setJobs] = useState<Job[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    axios
      .get<Job[]>('http://localhost:8080/api/jobs')
      .then(res => {setJobs(res.data); console.log(res.data)})
      .catch(err => {
        console.error(err);
        setError('Failed to fetch jobs');
      });
  }, []);

  return (
    <main style={{ padding: '2rem' }}>
      <h1>CI Job Dashboard</h1>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {!jobs && !error && <p>Loading...</p>}
      {jobs && jobs.length === 0 && <p>No jobs found.</p>}
      {jobs && jobs.length > 0 && <JobChart jobs={jobs} />}
    </main>
  );
}

export default App;
