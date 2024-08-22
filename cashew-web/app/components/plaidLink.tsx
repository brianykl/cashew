   
import { useEffect } from 'react';
import { usePlaidLink } from 'react-plaid-link';

interface PlaidLinkTokenProps {
  linkToken: string;
}


export default function PlaidLink({ linkToken }: PlaidLinkTokenProps) {
  const { open, ready } = usePlaidLink({
    token: linkToken,
    onSuccess: (public_token, metadata) => {
      console.log('Success:', public_token, metadata);
      // Handle success
    },
  });

  useEffect(() => {
    if (ready) {
      open();
    }
  }, [ready, open]);

  return null;
}