import { Component } from 'react';

import New from './New';

const options = [
  { value: 123, label: 'foo foo foo' },
  { value: 456, label: 'bar bar bar' },
  { value: 789, label: 'baz baz baz' },
];

export default class KitchenSink extends Component {
  render() { // eslint-disable-line
    return (
      <New
        isAdminOnly
        resourceName="kitchenSink"
        fields={{
          text: { name: 'Text' },
          email: { name: 'Email', type: 'email' },
          password: { name: 'Password', type: 'password' },
          bool: { name: 'Boolean', type: 'boolean' },
          date: { name: 'Date', type: 'date' },
          datetime: { name: 'Datetime', type: 'datetime' },
          select: {
            name: 'One thing',
            type: 'select',
            options,
          },
          multi: {
            name: 'Several things',
            type: 'multi',
            options,
          },
          textarea: { name: 'Textarea', type: 'textarea' },
        }}
      />
    );
  }
}
