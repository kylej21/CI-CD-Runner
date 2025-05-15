import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Tooltip,
  Legend
} from 'chart.js';

ChartJS.register(CategoryScale, LinearScale, BarElement, Tooltip, Legend);

type Job = {
  id: string;
  lint_errors: number;
  lint_warnings: number;
};

type Props = {
  jobs: Job[];
};

export default function JobChart({ jobs }: Props) {
  const data = {
    labels: jobs.map(job => job.id),
    datasets: [
      {
        label: 'Errors',
        data: jobs.map(job => job.lint_errors),
        backgroundColor: 'rgba(255, 99, 132, 0.7)',
      },
      {
        label: 'Warnings',
        data: jobs.map(job => job.lint_warnings),
        backgroundColor: 'rgba(255, 206, 86, 0.7)',
      },
    ],
  };

  return (
    <div>
      <h2>Lint Issues Per Job</h2>
      <Bar data={data} />
    </div>
  );
}
