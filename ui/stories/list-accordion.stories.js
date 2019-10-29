/* eslint-env node */
import hbs from 'htmlbars-inline-precompile';
import notes from './list-accordion.md';
import productMetadata from '../app/utils/styleguide/product-metadata';

export default {
  title: 'ListAccordion',
  parameters: {
    notes,
  },
};

export const Accordion = () => {
  return {
    template: hbs`
      <h5 class="title is-5">Accordion</h5>
      <ListAccordion @source={{products}} @key="name" as |ac|>
        <ac.head @buttonLabel="details">
          <div class="columns inline-definitions">
            <div class="column is-1">{{ac.item.name}}</div>
            <div class="column is-1">
              <span class="bumper-left badge is-light">{{ac.item.lang}}</span>
            </div>
          </div>
        </ac.head>
        <ac.body>
          <h1 class="title is-4">{{ac.item.name}}</h1>
          <p>{{ac.item.desc}}</p>
          <p><a href="{{ac.item.link}}">Learn more...</a></p>
        </ac.body>
      </ListAccordion>
          `,
    context: {
      products: productMetadata,
    },
  };
};

export const OneItem = () => {
  return {
    template: hbs`
      <h5 class="title is-5">Accordion, one item</h5>
      <ListAccordion @source={{take 1 products}} @key="name" as |a|>
        <a.head @buttonLabel="details">
          <div class="columns inline-definitions">
            <div class="column is-1">{{a.item.name}}</div>
            <div class="column is-1">
              <span class="bumper-left badge is-light">{{a.item.lang}}</span>
            </div>
          </div>
        </a.head>
        <a.body>
          <h1 class="title is-4">{{a.item.name}}</h1>
          <p>{{a.item.desc}}</p>
          <p><a href="{{a.item.link}}">Learn more...</a></p>
        </a.body>
      </ListAccordion>
          `,
    context: {
      products: productMetadata,
    },
  };
};

export const NotExpandable = () => {
  return {
    template: hbs`
      <h5 class="title is-5">Accordion, not expandable</h5>
      <ListAccordion @source={{products}} @key="name" as |a|>
        <a.head @buttonLabel="details" @isExpandable={{eq a.item.lang "golang"}}>
          <div class="columns inline-definitions">
            <div class="column is-1">{{a.item.name}}</div>
            <div class="column is-1">
              <span class="bumper-left badge is-light">{{a.item.lang}}</span>
            </div>
          </div>
        </a.head>
        <a.body>
          <h1 class="title is-4">{{a.item.name}}</h1>
          <p>{{a.item.desc}}</p>
          <p><a href="{{a.item.link}}">Learn more...</a></p>
        </a.body>
      </ListAccordion>
          `,
    context: {
      products: productMetadata,
    },
  };
};
