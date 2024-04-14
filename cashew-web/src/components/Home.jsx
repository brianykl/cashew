import logo from '../assets/cashew.svg';
import { Login } from '../components/Login';
import { useNavigate } from 'react-router-dom';

export function Home() {
	const navigate = useNavigate();

    const goToRegistration = () => {
        navigate('/register');
    };

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <a>
          cash doesn't grow on trees.
        </a>
        <Login />
        <a>
          <button onClick={goToRegistration}> user registration</button>
        </a>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          check out my github
        </a>
      </header>
    </div>
  );
}

