import { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { jwtDecode } from "jwt-decode";
import { GoogleLogin } from '@react-oauth/google';
import { useNavigate } from 'react-router-dom';

export function LoginPage() {
  let navigate = useNavigate();
  const [jwtToken, setJwtToken] = useState('');

  // load the jwttoken from localstorage
  useEffect(() => {
    const storedToken = localStorage.getItem('jwtToken');
    if (storedToken) {
      setJwtToken(storedToken);
    }
  }, []);

  // verify creds and add in db if not added already
  const handleLoginSuccess = async (credentialResponse) => {
    console.log("Login was successful!")
    const idToken = credentialResponse.credential;
    localStorage.setItem('jwtToken', idToken);
    setJwtToken(idToken);

    const response = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${jwtToken}`,
        'Content-Type': 'application/json'
      }
    });   

    if (response.ok) {
      const data = await response.json();
      console.log('Login status:', data);
      return navigate("/home");
    } else {
      console.log('Login request failed');
    }
  };

  const handleLoginError = () => {
    console.log('Login Failed');
  };

  // logout 
  const handleLogout = () => {
    localStorage.removeItem('jwtToken');
    setJwtToken('');
  };

  return (
    <div> 
      <div className="mx-auto max-w-sm space-y-8">
        <div className="space-y-2 text-center">
          <h1 className="text-3xl font-bold">Online Judge System</h1>
          <p className="text-gray-500 dark:text-gray-400">Sign in to continue</p>
        </div>
        <div className="space-y-4"> 
          <GoogleLogin
            onSuccess={handleLoginSuccess}
            onError={handleLoginError}
          />
        </div>
      </div>
    </div>
  ); 

  function ChromeIcon(props) {
    return (
      <svg
        {...props}
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <circle cx="12" cy="12" r="10" />
        <circle cx="12" cy="12" r="4" />
        <line x1="21.17" x2="12" y1="8" y2="8" />
        <line x1="3.95" x2="8.54" y1="6.06" y2="14" />
        <line x1="10.88" x2="15.46" y1="21.94" y2="14" />
      </svg>
    )
  }
}
