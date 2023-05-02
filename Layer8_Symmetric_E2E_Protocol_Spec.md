# Layer8 E2E Symmetric Encryption Protocol: Software Competition Specification and Criteria

## Summary

**Written**: April 2023

**Author(s)**: Ravi Seyed-Mahmoud

**Purpose**: This document summarizes the software specifications that outline the requirements for all entries to the Layer8 E2E Symmetric Encryption Protocol competition. 

**Background**: Layer8 is intended to be a suite protocols that implement standard cryptographic primitives and VPN techniques to facilitate end-to-end encryption services between a user-agent (e.g., a browser) and a Service Provider (e.g., any participating website or API). It distinguishes itself from other VPN implementations by running exclusively within the browser as a WASM module that connects to the intended end point Service Provider through the Layer8 Anonymizing Reverse Proxy. Currently Layer8 exists only as a proof of concept (see HTTPs://github.com/satsite13/Layer8). 
Layer8 development is supported by Globe&Citizen: A fully remote startup dedicated to becoming Web3.0’s premier, distributed, news service (see HTTPs://github.com/globe-and-citizen). Its members currently span 7 nations: Canada, Brazil, Togo, Nigeria, Uzbekistan, Pakistan, and Bangladesh.

**This Competition Specifically**: You will deliver a Golang WASM module that can be distributed via a CDN that successfully connects to a custom HTTPs server, also written by you, in Golang. All standard Golang libraries are allowed in addition to well established packages (e.g., gorilla/mux).
Future Competitions: This is the first of a series of competitions to implement the suite of protocols necessary to make Layer8 work as a beta release. Upcoming implementation competitions include the Layer8 Asymmetric Key Exchange Protocol; Layer8 JWT Authentication Protocol (based on OIDC); and the Layer8 Forwarding Protocol (list expected to change and modify with development). 
Regardless of your success in any one competition, you are strongly encouraged to continue participating in all subsequent competitions and to engage with the members of Globe&Citizen in an ongoing manner. Once Layer8 enters production, staff developers will be required and your participation during the competition phase—May through December 2023—will be your primary distinguishing criteria if you are interested in applying. 

**Timeline**: Beta release is planned for approximately January to March 2024 with the aforementioned coding / development competitions occurring from now until then.

**Licensing**: All code will be released under the GPL-2.0 license and posted to GitHub – both submitted code for competition and once in production. You must avoid using libraries that do not comply with the GPL-2.0 license.
 
**Glossary**:
•	Connection Database: An in memory store of active connections (represented by CINs) maintained by a Layer8 server. 
•	Connection Identification Number (CIN):  A unique identifier assigned by the central server to each client connected.  
•	Layer8 Anonymizing Reverse Proxy (L8ARP): The standalone HTTPs reverse proxy responsible for scrubbing all identifying information from client HTTPs messages and facilitating the connection between client and service provider. The L8ARP is responsible for end user identification and obfuscation of this information from the Service Provider. 
•	Layer8 Interpreter: A version set of implementation rules for decoding and understanding Data Bundles. Found on both the client module and the server implementation. 
•	Service Provider: What would traditionally be known simply as the “server” or “website” or “web app.” The end point that the client is attempting to connect to.

## Specifications & Description

### Introduction to the Layer8 E2E Symmetric Encryption Protocol
This specification outlines that portion of the Layer8 Protocol Suite that is responsible for transmission of symmetrically encrypted HTTPs messages between a client browser and Service Provider. It is to be fully implemented within the seventh layer of the OSI model (hence the name, “Layer8”). Its implementation is to be transparent to all lower layers of the network. In other words, symmetrically encrypted HTTPs messages should be indistinguishable from regular HTTPs messages having the content-type of text/plain and carrying a pseudorandom assortment of base64 ASCII characters in the request/response body. 
This is the first of several protocols, to be implemented. Other envisioned / proposed protocols include the Layer8 Asymmetric Key Exchange Protocol; Layer8 JWT Authentication Protocol; and the Layer8 Forwarding Protocol. In sum, these protocols together will deliver a transparent, browser based, end-to-end encrypted service that provides online anonymity over standard HTTPs enabled devices. Ultimately, the purpose is to de-associate a user’s true identity from their online content choices. In short, the goal is to create a situation where only the user is aware of both their true identity and content choices simultaneously while making it impossible for any other actor to have this information.
To accomplish this goal, the L8ARP will ultimately act as a CDN like distributor of vetted web apps that implement the Layer8 protocol adding a much needed security check to the World Wide Web (Kobeissi, 2021). An end user who authenticates with Layer8 should be assured that their true identity is dissociated fully, and irreversibly, from their online content choices.

![image](https://user-images.githubusercontent.com/116566901/235720960-8f589ffb-856c-4920-8aec-5280b6941538.png)

__Fig. 1: Highlevel Architecture of a future Layer8 System.__

## Task Overview
(For the purposes of this initial implementation / contest, secure transmission of symmetric keys is to be assumed. Predetermined symmetric keys are to be used with the assumption that an asymmetric key exchange protocol will be developed and deployed in the future.) 
You are to write:
1)	An in browser Golang WASM module, deliverable over CDN, that connects to…
2)	…a stand alone, HTTPs server, also written by you in Golang.
(Note: Fig. 1 above suggests 3 entities. This is a future vision. You will only be implementing your protocol between 2 entities: a client and a server.
The Golang WASM module functionality is to be incorporated into the global window object as L8 (i.e., window.L8.<methodCall> ). The WASM object, hereafter referred to as L8, is to expose the following methods:
1)	connect(serverUrl string, [options JSON])(connectionID uint16)
2)	sendEncrypted(content string, [options json])(r HTTP_Response)
The method connect(), at present, is simply a mock implementation of what will be the asymmetric key exchange between the client and server. Returned is the Connection Identification Cumber (CIN) assigned by the L8 Server (currently a 16 bit uint). 
The method L8.sendEncrypted() acts on the aforementioned connection to execute the following: 
1.	Encode the user’s chosen content as binary data.
2.	Chunk the content into segments. 
3.	Symmetrically encrypt each segment using an established cipher (e.g., AES-256)
4.	Package each segment into a Data Bundle (as described later).
5.	Populate the custom and standard HTTPs headers (as described later).
6.	Calculate a message authentication code using an established algorithm (e.g., HMAC) over the entire Data Bundle + Custom HTTPS headers. Append to the encrypted Data Bundle.
7.	Convert the Encrypted Data Bundle and MAC from binary to b64 encoded ASCII text.
8.	Send a standard HTTPs POST request composed of the custom HTTPs headers, Data Bundle, and MAC to the Layer8 server where it is to be reconstituted and interpreted.
Upon receipt of the POST request(s), your Layer8 Interpreter is to reverse the above process using the preestablished, dummy, symmetric key(s) “shared” with the client. Special attention should be made to de-segmenting the content for reconstitution by the server into its original form. 
General Considerations
-	The client side WASM module is to use standard browser APIs (i.e., fetch) so that all Layer8 functioning is transparent to lower networking levels (e.g., HTTPs, TLS, TCP, IP, etc.).
-	You are to make use of standard HTTPs message formatting so that extension to the protocol such as HTTPs/2 are utilized by default. 
-	Future implementations will make use of Websockets BUT DO NOT use Websockets in your current implementation. 
-	Assume that all symmetric keys have been securely and appropriately transmitted and stored. 
-	A future, in browser, key management solution will be implemented. For now, any storage location is appropriate (i.e., local storage or indexDB).
-	Cookies should be unnecessary and, in fact, are not allowed.

### Visual Overview of a Message With Encrypted Body
  
**HTTPs Headers** {
POST </server-api-endpoint> <HTTPs>
  Content-Type: Plain/Text
  Content-Length: <length> 
  < …other HTTPs headers… >

**Custom Headers** {
  x-cin: <16 bits (b64)>
  x-msg-cntr: <32-128 bits (b64)>

**Body** {
  Q29udHJhcnkgdG8gcG9wdWxhciBiZWxpZWYsIExvcmVtIElwc3VtIGlzIG5vdCBzaW1wbHkgcmFuZG9tIHRleHQuIEl0IGhhcyByb290cyBpbiBhIHBpZWNlIG9mIGNsYXNzaWNhbCBMYXRpbiBsaXRlcmF0dXJlIGZyb20gNDUgQkMsIG1ha2luZyBpdCBvdmVyIDIwMDAgeWVhcRoaXMgYm9vayBpcyBhIHRyZWF0aXNlIG9uIHRoZSB0aGVvcnkgb2YgZXRoaWNzLCB2ZXJ5IHBvcHVsYXIgZHVyaW5nIHRoZSBSZW5haXNzYW5jZS4gVGhlIGZpcnN0IGxpbmUgb2YgTG9yZW0gSXBzdW0sICJMb3JlbSBpcHN1bSBkb2xvciBzYgTG9yZW0YgTG9yZW0YgTG9yZW0bslbSBpcHN1bSBkb2xvciBzYgTG9yZW0YNsNsCCCWTG9yZWTG9y
  aXQgYW1ldC4uIIwgY29tZXMgZnJvbSBhIGxpbmUgaW4gc2VjdGlvbiAxLjEwLjMy


**Mac** {
    gYW1ldC4uIIwgY29tZXMgZnJvbSBhIGxpbmUgaW4gc2VjdGlvbiAxLjEwLjMy

_Fig. 2: Example HTTPs Message_

## Custom HTTPS Headers

### Connection Identification Number (CIN)
This is a (suggested) 16 bit number assigned by the server. It uniquely identifies every client connected and will be used to look up connection parameters and details during the creation, maintenance, and servicing of a connection. For the purposes of this implementation, the CIN can be requested via a very trivial request/response cycle invoked by L8.connect(). Clients without CIN yet assigned can set this field to zero and accept any subsequently assigned CIN. (In a subsequent competition, a real CIN assignment procedure is to be defined during development of the Layer8 Key Exchange protocol.)
CIN values must uniquely identify clients to the server such that your implementation can distinguish between all active connections. Implied by the 16 bit value is the ability of your server to maintain up to approximately 64,000 simultaneous, active, connections. (Note: 64,000 is not a hard requirement given that no hardware specifications are included here. Rather, it is an arbitrary value. The best implementation will be able to maintain the most connections relative to competing implementations.)
The server MUST maintain a database or in memory lookup table of active connections uniquely identified by the CIN. CINs are to be recycled when a connection is closed and reused as necessary. 
A variable subset of CIN numbers (e.g., 0 – 100) may be reserved for special usage.

### Data Bundle Format
(Note: the top level Data Bundle, that which includes the HTTPs custom headers is a special case. All other nested Data Bundles include their headers within the encrypted payload.)
“Data Bundle” refers to a set of headers + encrypted payload + associated MAC computed over both. All Data Bundles MUST start with an 8 bit header value that indicates the type of Data Bundle. Implied is the possibility of 256 different types of Data Bundle that could be defined in the future. 
Every data bundle always has the following headers followed by some combination of custom headers, defined by you in your implementation, and attached to a Data Bundle Type as follows:

0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ 
[Data Bundle Type| Length                        | Flags        ] 
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
[Maj V. |Min V. |
+-+-+-+-+-+-+-+-+
Fig. 3: Headers Standard to All Data Bundles

After your Layer8 Interpreter decrypts a Data Bundle, it reads the first 8 bits to determine the type of bundle which, through a standardized specification, indicates to the Layer8 Interpreter the number and type of header(s). Your interpreter will then read the next 16 bits to determine the bit length pertaining to the current Data Bundle (including any MAC). Eight flags, whose importance remains TBD, are then interpreted followed by the major / minor version.     

An example Data Bundle would be: 
0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ 
[Data Bundle Type| Length                        | Flags      ] 
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
[V. Maj |V. Min | Custom Header | Custom Header | Custom Header |                             
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
[ Custom Header | Custom Header | Custom Header | Custom Header ]
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
[ Payload +/- IV                                                ]
|                                                               |
~                                                               ~
~                                                               ~                                                               |                                                               |  
[                                                               ]  
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
[ Random, Variable Length, Padding                              ]
~                                                               ~
[                                                               ]
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
[ MAC (covering the padding, payload and headers)               ]
~                                                               ~
[                                                               ]
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  

_Fig. 4: An Example Data Bundle_

The purpose of this arrangement is to allow for rapid, recursive processing of data bundles that can be processed in parallel. 
The inspiration of the Data Bundles comes from the structure of JSON. Each payload is actually just an array of JSON objects. 
[{
  type: 1,
  length: 900,
  flags: 00010101,
  payload: [{
    type: 00000001,
    length: 900,
    flags: 00010101,
    payload: [{…}, {…}, {…}],
    mac: Y29taX14…uIIwgZQgYWXldCM
  },
  {
    type: 33,
    length: 96,
    flags: 00111101,
    customHeader: “”,
    customHeader: “”,
    customHeader: “”,
    payload: [{…}, {…}, {…}],
    mac: aXQgYW1ldC4…uIIwgY29tZXM
   },
   { …
   }],
  mac: G8gcG9wdWx9…tZXhciBiZWxp
}]
    
_Figure 3: Notice the three levels of payload, 1st (yellow), 2nd (green), 3rd (pink). _

In the above example, Type 33, is and invented Data Bundle that communicates to the Layer8 Interpreter what its headers are, what bytes correspond to its payload, and how long its total length is. This allows any Data Bundle (more precisely, those bytes that pertain to a Data Bundle) to be sliced from the ASCII blob and passed in their entirety to a service worker for processing.

#### Data Bundles are Recursive 

The terms Payload and Data Bundle are closely related but not equivalent.
-	A Data Bundle refers to the combination of Headers + Payload + MAC.
-	A Payload refers to only that portion of the Data Bundle that is NOT the MAC or Headers.
-	A Payload may itself be another bundle meaning that it too contains Headers + Payload + MAC. 
Payloads can hold one or more Data Bundles that may or may not be further encrypted. (Note the purpose of this complex system is to provide for future protocol extension). In fact, a payload is really just an array of one or more Data Bundles (i.e., JSON objects).
The top-level Data Bundle of any HTTPs message is special. It is always the combination of the three, custom, HTTPs headers (x-cin, x-msg-cntr), the b64 payload and the appended MAC. 

## Data Bundle Types
### Experimental, 0000000
This data bundle type must only be used during research and development. The only mandatory header are the generic headers of Fig. 3. +/- those custom headers you suggest. Any Layer8 server in production must automatically discard any data bundle of this type. 

### Recursive, 00000001
This is the primary data type and represents an object. It, itself, always begins with a single byte signifying the data bundle payload type and length.
    
### Dummy, 00000002
The data bundle type of 002 represents a dummy data bundle. It is used for obfuscation and confusion. Its headers are yet to be defined.

## Windowing, Chunking, and Reconstitution

Ultimately, the Layer8 system is envisioned to work over a full duplex connection such as Websockets. However, for the purposes of this implementation, you are to make use of HTTPs with the assumption that the server cannot, unilaterally, send messages. Fort the purposes of this competition, you are to assume that the client can only send requests and that the server can only send responses. (Theoretically, all modern browsers, under the hood, implement HTTP/2 transparently. This should, in theory, enable full duplex communication over the lifetime of an HTTPs request/response cycle. However, the goal is to make Layer8 transparent to the end user meaning that invoking L8.sendEncrypted() should behave near equivalent to the native fetch()and carry a near equivalent function signature.) 

In time, and once Layer8 is implemented over a full duplex connection, a bytewise TCP style, sliding, windowing algorithm is to be developed for data integrity. For the purposes of this competition, however, you are to implement a “Stop-n-Wait, Hail Mary” protocol as follows:
•	The client is to track, on a byte-by-byte level, all data sent and acknowledged.
•	The client maintains a SEND window which mirrors a RECEIVE window on the server. 
•	The server is to advertise, through the use of custom response headers, the state of it’s RECEIVE window.  
•	Bytes in the client’s SEND window are to be categorized as one of the following: sent & acknowledged; sent but not yet acknowledged; not sent but ready to be received; not sent and not ready to be received.
•	Bytes in the server’s RECEIVE window are to be categorized as one of the following: received and acknowledged; ready to be received; not ready to be received. 
•	If the client is required to send data spread across multiple messages, it is to communicate all necessary information regarding the total message to be received to the server using custom headers.
•	The server is to acknowledge all data received on a byte-by-byte basis. 
•	The client is to STOP sending as necessary.
•	The client is to WAIT for acknowledgements before resending messages presumed lost (i.e., “Stop-n-Wait”)
•	The server, by contrast, is to acknowledge the receipt of all messages and bytes and then, in a single streamed HTTPs message, respond on the assumption that the underlying internet protocols will work as expected. The server is then free to close the connection if appropriate (i.e., a “Hail Mary”, return to the client, of the requested resource with the assumption of success).

Notwithstanding the above, it is possible that advances in HTTPs/2 make the process of bytewise, high fidelity, reconstitution trivial and largely implemented by default. If so, this would represent a substantial advantage and you will not be penalized for making use of default functionality. In fact, achieving full functionality with a minimum of custom code is the ideal. However, if your chunking / windowing solution relies on in built functionality, you are to identify and explain how it works.

## Standardized Goal 

For the purposes of this implementation, your primary objective will be the transmission of the English, King James Bible to the server – streamed over a single HTTPs request (almost certainly divided under the hood by HTTP/2) where it is to be received and reconstituted by the server. Once fully received, the server is to stream back, via a single HTTPs response containing the French version of the King James Bible after following the same, analogous, steps 1 – 8 outlined in section Task Overview.
Streaming responses bitwise is a routine HTTPs practice and therefore should be trivial to implement (hence the idea of the “Hail Mary” return of a message). However, streaming requests is less common and therefore may pose challenges. If necessary, multiple requests may be sent. In such a scenario, your reconstitution algorithm at the server should take such a possibility into account as per the above suggested “Stop-n-Wait, Hail Mary” protocol.

## Language Requirements
Submissions will be written in Golang primarily with JS, and/or Type Script used as necessary. The browser module is to be compiled into WASM and the endpoint server is to be written in Golang. 
Where possible, only the latest, long term stable release of a software should be used.

## Criteria for Judging
The client module will be judged according to the request and response processing times. Network transmission times are to be ignored across the network. 
Your server implementation will be judged according to request and time to response processing times in addition to the number of simultaneous connections it can service.
Given that the messages are being carried via HTTPs, you may need to use various response/request timers to ensure receipt of all relevant data during chunking and reconstitution. This is not, de facto, unacceptable. Unnecessary delays, however, should be minimized. 

Please note: 
“Premature optimization is the root of all evil” – Donald Knuth.
Your implementation is expected to be “good enough” which is why all implementors will be paid for their entries and the payout for “winning” only a portion of the winner’s total compensation.  

##Negotiating Algorithms
Future versions of the Layer8 E2E Symmetric Encryption Protocol may include negotiation of different cryptographic algorithms configured in various ways. Your implementation of version 0, however, need only include a single set of predetermined algorithms and configurations chosen by you.
Message Counter: 16  to 64 Bits(?)
For every message sent, message counter must be incremented by 1 such that every message sent between two endpoints is uniquely identified by the combination of CIN and message counter. 
CIN should initialize to 0 at both transmitting ends of a connection. 
Upon reaching the message counter maximum, a connection is to be refreshed. The exact mechanisms for refreshing a connection is to be determined in a future protocol (i.e., the Layer8 Asymmetric Key Exchange Protocol). 
    
    
## Padding and IVs
If your chosen encryption algorithm requires an initialization vector, this is to be concatenated prior to the encrypted payload. 
If your chosen encryption algorithm requires a specific block size, you are to incorporate provisions for adding padding into your implementation. 
In this, version 0, there is no need to incorporate variable length padding for the purposes of obfuscating the true length of the data bundle proper. However, you are encouraged to include, as part of your pre processor, the ability to include a random number of bits which can be appended to the payload. The length of this padding should be included as a custom header in the Data Bundle. Ultimately, this functionality will be included as part of a future subcomponent responsible for preparation and processing of the data bundles. At present, it would be only a bonus. 
Message authentication codes should be computed to include the entire data bundle regardless up variable length padding used. In other words, the message authentication code is indifferent to the contents of the encrypted payload. 

## Logged Events
The inevitability of malicious attempts means that the successful implementation will have provisions for maintaining a log of suspicious or otherwise strange behaviour. Which events, and what associated details included are left as a matter of implementation. At present, a simple log file print out is adequate and no database integration is necessary. None but standard expected details need be included (i.e., date & time, CIN number, message counter, & server response).
Treatment of Specific HTTPs Headers
•	Content-Length: <content-length in bytes>
(Carries the message length info for parsing.)
•	Cache-Control: no-store
(For v0.0 at least, no responses should be cached.)
•	Connection:  x-cin, x-msg-cntr 
(Use of this header is to be determined.)
•	Set-Cookie: <…>
(Currently there are no plans to use the Set-Cookie header and you should NOT set cookies.)

## Contact: 
Author: Ravi Seyed-Mahmoud
Email: ravi.student.dev@gmail.com
github: https://github.com/stravid87 
References:
1.	Kobeissi, N. (2021) An Analysis of the ProtonMail Cryptographic Architecture. Retrieved from https://eprint.iacr.org/2018/1121.pdf
