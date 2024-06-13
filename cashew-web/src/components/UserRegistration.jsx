import React, { useState } from 'react';
import logo from '../assets/cashew.svg';
import { useNavigate } from 'react-router-dom';

export function UserRegistration() {
    const [userData, setUserData] = useState({
        email: '',
        name: '',
        password: '',
        confirmPassword: ''
    });
    
    const [successMessage, setSuccessMessage] = useState('');
    const navigate = useNavigate();

    const handleChange = (event) => {
        const { name, value } = event.target;
        setUserData(prevState => ({
            ...prevState,
            [name]: value
        }));
    
    };

    const handleSubmit = (event) => {
        console.log("registering user...")
        event.preventDefault();
        if (userData.password !== userData.confirmPassword) {
            alert('passwords do not match!');
            return; 
        }

        console.log('user data:', userData);
        fetch('http://localhost:4000/v1/user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: userData.email,
                name: userData.name,
                password: userData.password

            })
        })
        .then(data => {
            console.log('success:', data);
            setSuccessMessage('user was successfully created!');
            setUserData({
                email: '',
                name: '',
                password: '',
                confirmPassword: ''
            });
            // redirect or do somehting if successful
        })
        .catch((error) => {
            console.error("error: ", error);
            alert("error during user registration")
        })
    };

    return (
        <div>
            <h1>user registration</h1>
            {/* <img src={logo} alt="Logo" className="logo"/> */}
            {successMessage && <div className="success-message">{successMessage}</div>}
            <form onSubmit={handleSubmit}>
                <input
                    type="text"
                    name="email"
                    value={userData.email}
                    placeholder="email"
                    onChange={handleChange}
                    required
                />
                <input
                    type="text"
                    name="name"
                    value={userData.name}
                    placeholder="name"
                    onChange={handleChange}
                    required
                />
                <input
                    type="password"
                    name="password"
                    value={userData.password}
                    placeholder="password"
                    onChange={handleChange}
                    required
                />
                <input
                    type="password"
                    name="confirmPassword"
                    value={userData.confirmPassword}
                    placeholder="confirm password"
                    onChange={handleChange}
                    required
                />
                <button type="submit">register</button>
            </form>
            <button className="home-button" onClick={() => navigate('/')}>go to home page</button>
        </div>
    );

}
