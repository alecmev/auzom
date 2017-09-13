import { Component as ReactComponent } from 'react';
import classNames from 'classnames';

export default class Component extends ReactComponent {
  constructor(...args) {
    super(...args);
    // ugly but works; will be obsolete with CSS modules
    try {
      require(`../containers/${this.constructor.name}.scss`); // eslint-disable-line
    } catch (err1) {
      try {
        require(`../components/${this.constructor.name}.scss`); // eslint-disable-line
      } catch (err2) {
        // empty
      }
    }
  }

  cn({
    u: utilities,  // class="u-utility"
    d: descendant, // class="Component-descendant"
    m: modifiers,  // class="Component--modifier"
    s: states,     // class="Component is-state"
    f: force,      // class="u-utility Component"
  }) {
    const base = (
      this.constructor.name +
      (descendant ? `-${descendant}` : '')
    );

    const map = (a, prefix) => {
      if (!a) return false;
      return (Array.isArray(a) ? a : [a]).map(x => x && (prefix + x));
    };

    return classNames(
      map(utilities, 'u-'),
      (!utilities || descendant || modifiers || force) && base,
      map(modifiers, `${base}--`),
      map(states, 'is-'),
    );
  }

  cni(x = { }) {
    let f = x.f;
    if (x.u && x.f === undefined) f = true;
    return classNames(this.cn({ ...x, f }), this.props.className);
  }

  refCache = {}
  r = {}
  rcb = (name) => {
    if (!this.refCache[name]) {
      this.refCache[name] = (x) => {
        this.r[name] = x;
        return this.r[name];
      };
    }

    return this.refCache[name];
  }
}

// TODO: make sure the "utilities + base" heuristics don't have any undesired
// side-effects, and that the namespace is indeed redundant
