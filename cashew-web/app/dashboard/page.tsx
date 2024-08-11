'use client'
import { getSession } from "@auth0/nextjs-auth0";
import { redirect } from "next/navigation";
import {
    usePlaidLink,
    PlaidLinkOptions,
    PlaidLinkOnSuccess,
  } from 'react-plaid-link';

export default async function Dashboard() {
    const session = await getSession();

    if(!session) {
        redirect('api/auth/login')
    }

    // const config: PlaidLinkOptions = {
    //     onSuccess: (public_token, metadata) => {},
    //     onExit: (error, metadata) => {},
    //     onEvent: (eventName, metadata) => {},
    //     token: 'GENERATED_LINK_TOKEN'
    // }

    // const {open, exit, ready} = usePlaidLink(config)
    // // if (ready) {
    // //     open()
    // // }

    return (
      <div>cashew dashboard
        <a href="/api/auth/logout">Logout</a>
      </div>
    );
  }