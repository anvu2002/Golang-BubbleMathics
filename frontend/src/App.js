import React, { useEffect, useState } from 'react';

const App = () => {
    const [questions, setQuestions] = useState([]);
    const [messages, setMessages] = useState([]);

    useEffect(() => {
        // Fetch questions from the API
        fetch('/api/questions')
            .then(response => response.json())
            .then(data => setQuestions(data))
            .catch(error => console.error('Error fetching questions:', error));

        // Set up WebSocket connection
        const socket = new WebSocket('ws://localhost:8080/ws');

        socket.onopen = () => {
            console.log('WebSocket connection established');
        };

        socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            setMessages(prevMessages => [...prevMessages, message]);
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        socket.onclose = () => {
            console.log('WebSocket connection closed');
        };

        return () => {
            socket.close();
        };
    }, []);

    return (
        <div>
            <h1>Welcome to BubbleMathics</h1>
            <h2>Questions</h2>
            <ul>
                {questions.map((question, index) => (
                    <li key={index}>{question.text}</li>
                ))}
            </ul>
            <h2>Messages</h2>
            <ul>
                {messages.map((message, index) => (
                    <li key={index}>{JSON.stringify(message)}</li>
                ))}
            </ul>
        </div>
    );
};

export default App;
