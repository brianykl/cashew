import React from 'react';

export function Dashboard() {
	const email = sessionStorage.getItem('email');

	const handlePlaidLink = () => {
		// Implement Plaid Link integration here
		console.log('plaid link button clicked');
	};

	return (
		<div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
			<h1>welcome, {email}</h1>
			<button onClick={handlePlaidLink}>connect to bank account</button>
		</div>
	);
}