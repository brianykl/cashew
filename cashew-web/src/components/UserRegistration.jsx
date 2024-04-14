import React, { useState } from 'react';
import logo from '../assets/cashew.svg';
import { useNavigate } from 'react-router-dom';

export function UserRegistration() {
    const [userData, setUserData] = useState({
        email: '',
        password: '',
        confirmPassword: ''
    });

    const handleChange = (event) => {
        const { name, value } = event.target;
        setUserData(prevState => ({
            ...prevState,
            [name]: value
        }));
    
    };

    const handleSubmit = (event) => {
        event.preventDefault();
        if (userData.password !== userData.confirmPassword) {
            alert('passwords do not match!');
            return; 
        }
        console.log('user data:', userData);
        // registration logic, call user microservice
    };

    return (
        <div>
            <h1>user registration</h1>
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
        </div>
    );

}
