import { useState, FormEvent, ChangeEvent } from 'react';
import { useRouter } from 'next/router';
import InputMask from 'react-input-mask';

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
    // Assuming you will handle the form data here or pass it to the next page
    router.push('/categories'); // Update the path as per your flow
  };

  return (
    <div className="container">
      <h1>Enter Your Information</h1>
      <form onSubmit={handleSubmit}>
        <label htmlFor="name">Name:</label>
        <input
          type="text"
          id="name"
          value={name}
          onChange={handleNameChange}
          required
        />
        <label htmlFor="startDate">Budget Start Date:</label>
        <InputMask
          mask="9999-99-99"
          value={startDate}
          onChange={handleStartDateChange}
          placeholder="yyyy-mm-dd"
        >
          {(inputProps) => <input {...inputProps} type="text" />}
        </InputMask>
        <label htmlFor="endDate">Budget End Date:</label>
        <InputMask
          mask="9999-99-99"
          value={endDate}
          onChange={handleEndDateChange}
          placeholder="yyyy-mm-dd"
        >
          {(inputProps) => <input {...inputProps} type="text" />}
        </InputMask>
        <button type="submit">Next</button>
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