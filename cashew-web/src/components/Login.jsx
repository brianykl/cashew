import React, { useState } from 'react';

export function Login() {
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');

	const handleSubmit = (event) => {
		event.preventDefault();
		console.log('login with:', email, password);
		// api call to authentication service
		fetch('http://localhost:4000/v1/user/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: email,
                password: password

            })
        })
		.then(response => response.json()) 
        .then(data => {
			console.log('this is data', data)
            if (data.token) {
			console.log('success:', data);
			sessionStorage.setItem('token', data.token);
			sessionStorage.setItem('email', data.user);
			console.log('JWT stored successfully:', data.token);
            // redirect or do somehting if successful
			window.location.href = '/dashboard';
			} else {
				alert('unsuccessful login')
			}
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
