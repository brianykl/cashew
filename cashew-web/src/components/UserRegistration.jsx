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
        fetch('http://localhost:3000/v1/user/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: userData.email,
                name: userData.name,
                password: userData

            })
        })
        .then(data => {
            console.log('success:', data);
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
                    oonCHange={handleChange}
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
