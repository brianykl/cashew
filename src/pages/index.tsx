import { useRouter } from 'next/router';
import { FormEvent } from 'react';

const Home: React.FC = () => {
  const router = useRouter();

  return (
    <div className="container">
      <img src="/logo.png" alt="cashew logo" className="logo" />
      <h1> welcome to cashew </h1>
      <p> a simple tool to help you make a pretty fancy budget </p>
      <button onClick={() => router.push('name-and-period')}>
        get started 
      </button>
      

      <style jsx>{`
        .container {
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
          height: 100vh;
          text-align: center;
          background-color: #e0d5c6; /* Light beige color */
          padding: 20px;
        }
        .logo {
          width: 500px;
          height: auto;
          margin-bottom: -200px;
        }
        h1 {
          font-size: 2.5rem;
          color: #333;
          margin-bottom: 10px;
        }
        p {
          font-size: 1.2rem;
          color: #666;
          margin-bottom: 20px;
        }
        button {
          padding: 10px 20px;
          font-size: 1rem;
          background-color: #6b8e23; /* Green color */
          color: #fff;
          border: none;
          border-radius: 4px;
          cursor: pointer;
        }
        button:hover {
          background-color: #556b2f; /* Darker green for hover effect */
        }
      `}</style>
    </div>
  );
};

export default Home;