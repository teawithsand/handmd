#!/usr/bin/env node

/**
 * Simple node script, which creates fuse.js index from data coming from stdin.
 */

const fuse = require("fuse.js")
const fs = require("fs")


const data = fs.readFileSync(0, 'utf-8');
const parsed = JSON.parse(data)

const fields = process.argv.slice(2)
const index = fuse.createIndex(fields, parsed)

process.stdout.write(JSON.stringify(index))