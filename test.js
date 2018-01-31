import foo from './foo';

function foo(bar) {
  console.log(bar);
  {
    let baz = 232;
  }
  return bar + 'kek';
}

foo.default = bar;