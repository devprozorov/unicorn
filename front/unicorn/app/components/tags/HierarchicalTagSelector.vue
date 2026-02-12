<template>
  <div class="hierarchical-tag-selector">
    <!-- Поиск по тегам -->
    <div class="search-section" v-if="enableSearch">
      <input
        v-model="localSearchQuery"
        type="text"
        :placeholder="t('tags.searchPlaceholder')"
        class="tag-search-input"
      />
    </div>

    <!-- Результаты поиска -->
    <div v-if="localSearchQuery && searchResults.length > 0" class="search-results-section">
      <div class="section-title">{{ t('tags.searchResults') }}</div>
      <div class="tags-grid">
        <label
          v-for="tag in searchResults"
          :key="tag"
          class="tag-checkbox-item"
        >
          <input
            type="checkbox"
            :checked="isSelected(tag)"
            @change="$emit('toggle-tag', tag)"
          />
          <span>{{ getLabel(tag) }}</span>
        </label>
      </div>
    </div>

    <!-- Категории как аккордеон -->
    <div v-else class="categories-accordion">
      <div
        v-for="category in allCategories"
        :key="category.value"
        class="category-section"
      >
        <button
          @click="toggleCategory(category.value)"
          class="category-header"
          :class="{ active: expandedCategory === category.value }"
        >
          <div class="category-info">
            <div class="category-name">{{ category.label }}</div>
          </div>
          <span class="expand-icon">{{ expandedCategory === category.value ? '-' : '+' }}</span>
        </button>

        <!-- Специализации внутри категории (вложенная структура) -->
        <div v-if="expandedCategory === category.value" class="category-content">
          <div
            v-for="spec in getSpecializationsForCategory(category.value)"
            :key="spec.value"
            class="spec-section"
          >
            <!-- Если есть дети - показываем как подраздел -->
            <div v-if="spec.children && spec.children.length > 0">
              <button
                @click="toggleSubcategory(spec.value)"
                class="subcategory-header"
                :class="{ active: expandedSubcategory === spec.value }"
              >
                <span>{{ spec.label }}</span>
                <span class="expand-icon-small">{{ expandedSubcategory === spec.value ? '-' : '+' }}</span>
              </button>

              <!-- Детали внутри специализации -->
              <div v-if="expandedSubcategory === spec.value" class="subcategory-content">
                <div class="tags-grid">
                  <label
                    v-for="detail in spec.children"
                    :key="detail.value"
                    class="tag-checkbox-item"
                  >
                    <input
                      type="checkbox"
                      :checked="isSelected(detail.value)"
                      @change="$emit('toggle-tag', detail.value)"
                    />
                    <span>{{ detail.label }}</span>
                  </label>
                </div>
              </div>
            </div>

            <!-- Если нет детей - показываем как чекбокс -->
            <label v-else class="tag-checkbox-item">
              <input
                type="checkbox"
                :checked="isSelected(spec.value)"
                @change="$emit('toggle-tag', spec.value)"
              />
              <span>{{ spec.label }}</span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- Выбранные теги внизу -->
    <div v-if="selectedTags.length > 0" class="selected-tags-section">
      <div class="section-header">
        <span>{{ t('tags.selected') }}: {{ selectedTags.length }}</span>
        <button @click="$emit('clear-all')" class="clear-btn">{{ t('tags.clearAll') }}</button>
      </div>
      <div class="selected-tags-list">
        <button
          v-for="tag in selectedTags"
          :key="tag"
          @click="$emit('toggle-tag', tag)"
          class="selected-tag"
        >
          {{ getLabel(tag) }}
          <span class="remove-icon">x</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useTags } from '~/composables/useTags'
import { useI18n } from '~/composables/useI18n'
import { tagLabelMap, tagTaxonomy, type TagHierarchy } from '~/utils/tags'

const props = defineProps<{
  selectedTags: string[]
  selectedCategory?: string | null
  selectedSpecialization?: string | null
  enableSearch?: boolean
  showRecommendations?: boolean
}>()

const emit = defineEmits<{
  'toggle-tag': [tag: string]
  'select-category': [category: string]
  'select-specialization': [specialization: string]
  'clear-all': []
}>()

const tags = useTags()
const { t } = useI18n()
const localSearchQuery = ref('')
const expandedCategory = ref<string | null>(null)
const expandedSubcategory = ref<string | null>(null)

// Categories in current language
const allCategories = computed(() => {
  // Показываем только 4 основные категории: IT, GameDev, Startup, Other
  const mainCategories = ['IT', 'GameDev', 'Startup', 'Other']
  
  return tagTaxonomy
    .filter(cat => mainCategories.includes(cat.value))
    .map(cat => ({
      value: cat.value,
      label: t(cat.labelKey)
    }))
})

// Computed properties
const searchResults = computed(() => {
  if (!localSearchQuery.value || localSearchQuery.value.length < 2) return []
  return tags.searchResults.value
})

// Methods
function isSelected(tag: string): boolean {
  return props.selectedTags.includes(tag)
}

function toggleCategory(categoryValue: string) {
  if (expandedCategory.value === categoryValue) {
    expandedCategory.value = null
    expandedSubcategory.value = null
  } else {
    expandedCategory.value = categoryValue
    expandedSubcategory.value = null
    emit('select-category', categoryValue)
  }
}

function toggleSubcategory(subcategoryValue: string) {
  if (expandedSubcategory.value === subcategoryValue) {
    expandedSubcategory.value = null
  } else {
    expandedSubcategory.value = subcategoryValue
  }
}

function getSpecializationsForCategory(categoryValue: string) {
  const category = tagTaxonomy.find(c => c.value === categoryValue)
  if (!category?.children) return []

  return category.children.map(spec => ({
    value: spec.value,
    label: getLabel(spec.value),
    children: spec.children ? spec.children.map(detail => ({
      value: detail.value,
      label: getLabel(detail.value)
    })) : undefined
  }))
}

function findLabelKey(tag: string, nodes: TagHierarchy[]): string | null {
  for (const node of nodes) {
    if (node.value === tag) return node.labelKey
    if (node.children) {
      const childKey = findLabelKey(tag, node.children)
      if (childKey) return childKey
    }
  }
  return null
}

// Universal tag label resolver.
function getLabel(tag: string): string {
  const labelKey = findLabelKey(tag, tagTaxonomy) ?? tagLabelMap[tag]
  if (!labelKey) return tag
  const translated = t(labelKey)
  return translated === labelKey ? tag : translated
}
// Watch for external changes
watch(localSearchQuery, (newVal) => {
  tags.searchQuery.value = newVal
})
</script>

<style scoped>
.hierarchical-tag-selector {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

/* Поиск */
.search-section {
  margin-bottom: 0.5rem;
}

.tag-search-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 16px;
  background: white;
  transition: border-color 0.2s;
}

.tag-search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

/* Результаты поиска */
.search-results-section {
  background: white;
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid #e5e7eb;
}

.section-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: #6b7280;
  margin-bottom: 0.75rem;
}

.tags-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 0.5rem;
}

.tag-checkbox-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  background: white;
}

.tag-checkbox-item:hover {
  background: #f9fafb;
  border-color: #d1d5db;
}

.tag-checkbox-item input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: #3b82f6;
}

.tag-checkbox-item span {
  font-size: 0.875rem;
  color: #374151;
}

/* Категории как аккордеон */
.categories-accordion {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.category-section {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.category-header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  background: white;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
}

.category-header:hover {
  background: #f9fafb;
}

.category-header.active {
  background: #eff6ff;
  border-bottom: 1px solid #e5e7eb;
}

.category-info {
  flex: 1;
}

.category-name {
  font-size: 1rem;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 0.25rem;
}

.expand-icon {
  font-size: 0.875rem;
  color: #9ca3af;
  transition: transform 0.2s;
}

.category-header.active .expand-icon {
  color: #3b82f6;
}

.category-content {
  padding: 1rem 1.25rem;
  background: #f9fafb;
  border-top: 1px solid #e5e7eb;
}

.spec-section {
  margin-bottom: 0.75rem;
}

.spec-section:last-child {
  margin-bottom: 0;
}

.subcategory-header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.9rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.5rem;
}

.subcategory-header:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
}

.subcategory-header.active {
  background: #eff6ff;
  border-color: #93c5fd;
  color: #1e40af;
}

.expand-icon-small {
  font-size: 0.75rem;
  color: #9ca3af;
  transition: transform 0.2s;
}

.subcategory-header.active .expand-icon-small {
  color: #3b82f6;
}

.subcategory-content {
  padding-left: 1rem;
  margin-bottom: 0.5rem;
}

/* Выбранные теги */
.selected-tags-section {
  padding: 1rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  font-size: 0.9rem;
  font-weight: 500;
  color: #6b7280;
}

.clear-btn {
  padding: 0.375rem 0.75rem;
  background: transparent;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 0.875rem;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-btn:hover {
  background: #f9fafb;
  border-color: #d1d5db;
}

.selected-tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.selected-tag {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 6px;
  font-size: 0.875rem;
  color: #1e40af;
  cursor: pointer;
  transition: all 0.2s;
}

.selected-tag:hover {
  background: #dbeafe;
  border-color: #93c5fd;
}

.remove-icon {
  font-size: 1.25rem;
  line-height: 1;
  color: #60a5fa;
}

/* Адаптивность для мобильных */
@media (max-width: 768px) {
  .tags-grid {
    grid-template-columns: 1fr;
  }

  .tag-checkbox-item {
    min-height: 48px;
  }
}

/* Для touch устройств */
@media (hover: none) and (pointer: coarse) {
  .category-header,
  .tag-checkbox-item {
    min-height: 48px;
  }

  .tag-checkbox-item input[type="checkbox"] {
    width: 24px;
    height: 24px;
  }
}
</style>
