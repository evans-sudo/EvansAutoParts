$(function() {
  var categorySelect = $('#category');
  var subcategorySelect = $('#subcategory');
  var typeSelect = $('#type');
  var resultContainer = $('#resultContainer');

  var modelId = $('#model').val();
  var yearId = $('#year').val();
  var engineId = $('#engine').val();

  var categories = window.categories;

  var activeCategory = null;
  var activeSubcategory = null;
  var activeType = null;

  categories.forEach(function(c) {
    categorySelect.append('<option value="' + c.id + '">' + c.name + '</option>');
  });

  function populateSubcategorySelect() {
    subcategorySelect.empty();
    activeCategory.subcategories.forEach(function(s) {
      subcategorySelect.append('<option value="' + s.id + '">' + s.name + '</option>');
    });
  }

  function populateTypeSelect() {
    typeSelect.empty();
    activeSubcategory.types.forEach(function(t) {
      typeSelect.append('<option value="' + t.id + '">' + t.name + '</option>');
    });
  }

  categorySelect.change(function() {
      for (var i = 0; i < categories.length; i++) {
        if (categories[i].id == this.value) {
            activeCategory = categories[i];
            populateSubcategorySelect();
            break;
        }
      }
  });

  subcategorySelect.change(function() {
      for (var i = 0; i < activeCategory.subcategories.length; i++) {
        if (activeCategory.subcategories[i].id == this.value) {
            activeSubcategory = activeCategory.subcategories[i];
            populateTypeSelect();
            break;
        }
      }
  });

  typeSelect.change(function() {
    var value = this.value;
    resultContainer.load("/parts?modelId=" + modelId +
      "&yearId=" + yearId +
      "&engineId=" + engineId +
      "&typeId=" + value +
      "&employeeNumber=" + $("#employeeNumber").val());
  });
});
