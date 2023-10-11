# Layer8
A suite of  network protocol implementations that sum to create an anonymizing reverse proxy dedicated to dissociating a user's true identity from their online content choices.  

## Why
Thus far, the onus is overwhelmingly placed on end users of the internet to achieve anonymity online which limits scalability (think using Tor, installing a VPN, etc.). There are, however, niche circumstances where user anonymity is desirable and / or necessary. Unfortunately, frictionless solutions using the browser as a platform for end-to-end encryption (think Proton Mail) are easily compromised because of problems associated with trust in the standard client / server model . A market opportunity is available, thus, in the online microservice ecosystem to provide frictionless anonymity services to a company’s end users. In addition to exploring this opportunity, Layer 8 also serves simultaneously as an R&D foundation for future MAIC projects.

## Key Performance Indicator
Realize a complex production system using modern, but already available, technologies to produce a novel microservice. In other words, successfully applied secondary, but not primary, research and development. For example, greenfield implementations of vetted cryptographic primitives in web assembly language.    

## What
Layer 8 is designed to be a scalable internet service platform that enables end-to-end encryption via the browser. This, in turn, enables a user’s true identity to be stripped from their content choices. Very broadly, the proof-of-concept, works as follows: 
1.	A content delivery network serves an in-browser module which exposes the L8 global object. 
2.	The developer invokes client side methods to build an encrypted tunnel to their backend through the Layer 8 reverse proxy ( e.g., L8.registerCitizenship(…) ).
3.	Through an algorithm inspired by OAuth2.0, Layer 8 establishes an encrypted tunnel using JSON Web Tokens. 
4.	By acting as an HTTP reverse proxy, Layer 8 strips identifying header information requests and replaces it with custom metadata suitable for public aggregation.
5.	Metadata can be collated by the Service Provider without fear of deanonymizing their users. 
6.	Because an encrypted tunnel has been established, Layer 8 is ignorant of a user’s content choices whereas the Service Provider is ignorant of the user’s true identity. Only the end user, according to the scheme proposed, is aware of both their true identity and their content choices. 

To succeed, Layer 8 must solve the fundamental problem(s) associated with trusting an unknown server to deliver an uncompromised application to the client with every new page load (see Kobeissi, N. (2021) An Analysis of the ProtonMail Cryptographic Architecture. Retrieved from https://eprint.iacr.org/2018/1121.pdf). To users, it will provide a free, anonymizing, authentication service. 

Layer 8 must be open source so that it can be vetted by the larger internet security community. Crowd scrutiny will be fundamental to establishing its brand identity as a viable, trusted, third party. It should be noted that, technologically, it is already possible to build Layer 8. In fact, the working proof of concept that I wrote using Type Script is available on Github (github.com/satsite13/Layer8.git). Risky primary research and development should not be necessary to realize the project and will be avoided.

## Usage

### Local Development
> For local development, the proxy is available as a docker container. First, ensure that you have [Docker](https://docs.docker.com/get-docker/) installed. Then, run the following command to pull and start the proxy:
>
> ```bash
> docker run -d --network host deestarks/layer8-proxy
> ```

1. In your frontend application, add the following script tag to the head of your index.html file:
    ```html
    <script type="module" src="https://auth-service-5mcaj.ondigitalocean.app/assets/cdn/interceptor/index.local.js"></script>
    ```
2. In your backend application, install the following npm package:
    > Note: Node version 20 or higher is recommended.
    ```bash
    npm i layer8-middleware
    ```
3. Also, in your backend application, add the following code to your application's entry point:
    ```javascript
    // import the layer8 middleware
    const layer8 = require('layer8-middleware');
    // use the layer8 middleware
    app.use(layer8);
    ```

    An example usage in an express application is shown below:
    ```javascript
    const express = require('express');
    const layer8 = require('layer8-middleware');

    const app = express();

    app.use(express.json());
    app.use(layer8);
    ```

### Production

1. In your frontend application, add the following script tag to the head of your index.html file:
    ```html
    <script type="module" src="https://auth-service-5mcaj.ondigitalocean.app/assets/cdn/interceptor/index.js"></script>
    ```
2. In your backend application, install the following npm package:
    ```bash
    npm i layer8-middleware
    ```
3. Also, in your backend application, add the following code to your application's entry point:
    ```javascript
    // import the layer8 middleware
    const layer8 = require('layer8-middleware');
    // use the layer8 middleware
    app.use(layer8);
    ```
    An example usage in an express application is shown below:
    ```javascript
    const express = require('express');
    const layer8 = require('layer8-middleware');

    const app = express();

    app.use(express.json());
    app.use(layer8);
    ```
