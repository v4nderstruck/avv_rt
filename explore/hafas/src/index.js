import { createClient } from "hafas-client";
import { profile } from 'hafas-client/p/avv/index.js';
const userAgent = 'BeepBoop';
const client = createClient(profile, userAgent);

const departures = await client.departures("000000001001");
console.log(JSON.stringify(departures, null, 2))
