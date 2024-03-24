import https, { type RequestOptions } from 'node:https';
import {parseDocument, DomUtils, ElementType} from 'htmlparser2';

const options: RequestOptions = {
  host: 'shareprices.com',
  path: '/indices/ftse-all-share/'
}

function parseTable(data: string) {
  // console.log(data);
  const dom = parseDocument(data)

  // console.log(dom);

  const table = DomUtils.findOne((el) => el.type === 'tag' && el.name === 'table' && el.attribs.id === 'tblConstituents', dom.childNodes, true)
  const tBody = DomUtils.findOne((el) => el.type === 'tag' && el.name === 'tbody', table?.childNodes || [], true)

  if (!tBody) {
    return
  }

  // const rows = DomUtils.find((el) => el.type === 'tag' && el.name === 'a' && el.attribs.class === '---home__table-display-name', tBody.childNodes, true, 2000)
    // @ts-expect-error - it IS elements
  const rows: Element[] = tBody.childNodes
    .filter(Boolean)
    .map((node) => node.type === 'tag' && DomUtils.findOne((el) => el.type === 'tag' && el.name === 'a' && el.attribs.class === '---home__table-display-name', node.childNodes, true))
    // .filter((node) => node && node.type === 'tag')

  // const keys = rows.map(row => row.)

  // @ts-expect-error - it IS elements
  console.log(JSON.stringify(rows.map(row => row?.firstChild?.data).filter(Boolean)));

  // const rows = DomUtils.find((el) => el.type === 'tag' && el.name === 'tr', table.children, true, 500)
  // // console.log(table.childNodes.map((node) =>  node.type === 'tag' && node.name));
  // console.log(rows, rows.length);
}

const req = https.request(options, (res) => {
  let data = ''
  res.on('data', (chunk) => {
    data += chunk;
  })

  res.on('end', () => {
    parseTable(data)
  })
})

req.end()


