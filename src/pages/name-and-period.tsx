import { useState, ChangeEvent, FormEvent } from 'react';
import { useRouter } from 'next/router';

const NameAndPeriod: React.FC = () => {
    const [name, setName] = useState<string>('');
    const [startDate, setStartDate] = useState<string>('');
    const [endDate, setEndDate] = useState<string>('');
    const router = useRouter();

    const handleNameChange = (e: ChangeEvent<HTMLInputElement>) => {
        setName(e.target.value);
      };
    
      const handleStartDateChange = (e: ChangeEvent<HTMLInputElement>) => {
        setStartDate(e.target.value);
      };
    
      const handleEndDateChange = (e: ChangeEvent<HTMLInputElement>) => {
        setEndDate(e.target.value);
      };
    
    const handleSubmit = (e: FormEvent) => {
        e.preventDefault();
        // something, assume we handle form data here or pass to next page
        router.push('/categories'); // next page
    };

    return (
        <div className="container">
            <h1>basic info</h1>
            <form onSubmit={handleSubmit}>
                <label htmlFor='name'> what's your name?</label>
                <input
                    type="text"
                    id="name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                />
                <label htmlFor='startDate'> when do you want this budget to start?</label>
                <input
                    type="date"
                    id="startDate"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    required
                />
                <label htmlFor='endDate'> when do you want this budget to end?</label>
                <input
                    type="date"
                    id="endDate"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                    required
                />
                <button type="submit">next</button>
            </form>
            <style jsx>{`
        .container {
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
          height: 100vh;
          text-align: center;
          background-color: #e0d5c6; /* Darker beige color */
          padding: 20px;
        }
        h1 {
          font-size: 2.5rem;
          color: #333;
          margin-bottom: 20px;
        }
        label {
          font-size: 1.2rem;
          margin-bottom: 10px;
          color: #333;
        }
        input {
          padding: 10px;
          font-size: 1rem;
          margin-bottom: 20px;
          border: 1px solid #ccc;
          border-radius: 4px;
          width: 300px;
          max-width: 80%;
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

export default NameAndPeriod;