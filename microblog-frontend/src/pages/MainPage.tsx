import React from 'react';
import * as Chakra from '@chakra-ui/react';

export default function MainPage() {
  return (
    <Chakra.Grid
      templateAreas={`"header header"
                  "nav main"
                  "nav footer"`}
      gridTemplateRows={'50px 1fr 30px'}
      gridTemplateColumns={'150px 1fr'}
      h='100vh'
      gap='1'
    >
      <Chakra.GridItem
        pl='2'
        bg='orange.300'
        area={'header'}
      >
        Header
      </Chakra.GridItem>
      <Chakra.GridItem
        pl='2'
        bg='pink.300'
        area={'nav'}
      >
        Nav
      </Chakra.GridItem>
      <Chakra.GridItem
        pl='2'
        bg='green.300'
        area={'main'}
      >
        Main
      </Chakra.GridItem>
      <Chakra.GridItem
        pl='2'
        bg='blue.300'
        area={'footer'}
      >
        Footer
      </Chakra.GridItem>
    </Chakra.Grid>
  );
}
