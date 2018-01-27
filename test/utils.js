function unorderedList() {
  return {
    items: [],
    freeSpaces: [],
    add(value) {
      let index = this.items.length;
      if (this.freeSpaces.length) {
        index = this.freeSpaces.pop();
      }
      this.items[index] = value;
      return index;
    },
    remove(index) {
      this.freeSpaces.push(index);
      this.items[index] = unorderedList.REMOVED_ITEM;
    }
  }; 
}
unorderedList.REMOVED_ITEM = Symbol('REMOVED_ITEM');

export default {
  unorderedList
};
