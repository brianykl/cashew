import React, { useState } from 'react';

export function Login() {
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');

	const handleSubmit = (event) => {
		event.preventDefault();
		console.log('login with:', email, password);
		// api call to authentication service
		fetch('http://localhost:3000/v1/user/', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: userData.email,
                password: userData

            })
        })
        .then(data => {
            console.log('success:', data);
            // redirect or do somehting if successful
        })
        .catch((error) => {
            console.error("error: ", error);
            alert("unsuccessful login attempt")
        })
	};

	const handleEmailChange = (event) => {
		setEmail(event.target.value);
	};

	const handlePasswordChange = (event) => {
		setPassword(event.target.value);
		
	};


	return (
		<form onSubmit={handleSubmit}>
		<label htmlFor="email">email:</label>
		<input
			type="email"
			id="email"
			value={email}
			onChange={handleEmailChange}
			required
		/>
		<label htmlFor="password">password:</label>
		<input
			type="password"
			id="password"
			value={password}
			onChange={handlePasswordChange}
			required
		/>
		<button type="submit">login</button>
	</form>
	);
}
